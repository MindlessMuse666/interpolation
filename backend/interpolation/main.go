package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/MindlessMuse666/interpolation/backend/core/interpolation"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type Config struct {
	Interpolation struct {
		Port        int    `toml:"port"`
		RedisAddr   string `toml:"redis_addr"`
		RabbitMQURL string `toml:"rabbitmq_url"`
		CacheTTL    int    `toml:"cache_ttl"`
	} `toml:"interpolation"`
}

type InterpolateRequest struct {
	Method  string                `json:"method" binding:"required"`
	Points  []interpolation.Point `json:"points" binding:"required,min=2"`
	TargetX float64               `json:"target_x" binding:"required"`
}

type InterpolateResponse struct {
	Method string                `json:"method"`
	Result float64               `json:"result"`
	Curve  []interpolation.Point `json:"curve"`
	Cached bool                  `json:"cached"`
}

type CalculationEvent struct {
	Method    string                `json:"method"`
	Points    []interpolation.Point `json:"points"`
	TargetX   float64               `json:"target_x"`
	Result    float64               `json:"result"`
	Timestamp time.Time             `json:"timestamp"`
}

var (
	cfg         Config
	redisClient *redis.Client
	amqpConn    *amqp.Connection
	amqpChan    *amqp.Channel
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if _, err := toml.DecodeFile("../../config/config.toml", &cfg); err != nil {
		// Try local path if not found (for different environments)
		if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
			log.Warn().Err(err).Msg("Could not load config.toml, using defaults")
			cfg.Interpolation.Port = 8081
			cfg.Interpolation.RedisAddr = "localhost:6379"
			cfg.Interpolation.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
			cfg.Interpolation.CacheTTL = 3600
		}
	}
}

func setupRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: cfg.Interpolation.RedisAddr,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Error().Err(err).Msg("Redis connection failed")
		redisClient = nil // Graceful degradation: continue without Redis
	} else {
		log.Info().Msg("Connected to Redis")
	}
}

func setupRabbitMQ() {
	var err error
	// Retry connection to RabbitMQ
	for i := 0; i < 10; i++ {
		amqpConn, err = amqp.Dial(cfg.Interpolation.RabbitMQURL)
		if err == nil {
			break
		}
		log.Warn().Err(err).Msgf("RabbitMQ connection failed, retrying in 5s... (%d/10)", i+1)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Error().Err(err).Msg("RabbitMQ connection failed after retries")
		return
	}

	amqpChan, err = amqpConn.Channel()
	if err != nil {
		log.Error().Err(err).Msg("Failed to open RabbitMQ channel")
		return
	}

	err = amqpChan.ExchangeDeclare(
		"calculations", // name
		"direct",       // type
		true,           // durable
		false,          // auto-deleted
		false,          // internal
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to declare exchange")
		return
	}
	log.Info().Msg("Connected to RabbitMQ")
}

func getCacheKey(req InterpolateRequest) string {
	data, _ := json.Marshal(req)
	hash := md5.Sum(data)
	return fmt.Sprintf("interpolation:%x", hash)
}

func publishEvent(req InterpolateRequest, result float64) {
	if amqpChan == nil {
		return
	}

	event := CalculationEvent{
		Method:    req.Method,
		Points:    req.Points,
		TargetX:   req.TargetX,
		Result:    result,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal event")
		return
	}

	err = amqpChan.Publish(
		"calculations",          // exchange
		"calculation.completed", // routing key
		false,                   // mandatory
		false,                   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish event")
	}
}

// @Summary Интерполяция данных
// @Description Вычисляет значение y для заданного x, используя выбранный метод интерполяции, и возвращает точки для графика.
// @Accept json
// @Produce json
// @Param request body InterpolateRequest true "Параметры интерполяции"
// @Success 200 {object} InterpolateResponse
// @Failure 400 {object} map[string]string "Неверный формат данных или неизвестный метод"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /interpolate [post]
func handleInterpolate(c *gin.Context) {
	var req InterpolateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Check Cache
	cacheKey := getCacheKey(req)
	if redisClient != nil {
		val, err := redisClient.Get(context.Background(), cacheKey).Result()
		if err == nil {
			var resp InterpolateResponse
			if err := json.Unmarshal([]byte(val), &resp); err == nil {
				resp.Cached = true
				c.JSON(http.StatusOK, resp)
				return
			}
		}
	}

	// 2. Calculate
	var resY float64
	var curve []interpolation.Point
	var err error

	switch req.Method {
	case "linear":
		resY, curve, err = interpolation.LinearInterpolation(req.Points, req.TargetX)
	case "lagrange":
		resY, curve, err = interpolation.LagrangeInterpolation(req.Points, req.TargetX)
	case "newton":
		resY, curve, err = interpolation.NewtonInterpolation(req.Points, req.TargetX)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "неизвестный метод интерполяции"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := InterpolateResponse{
		Method: req.Method,
		Result: resY,
		Curve:  curve,
		Cached: false,
	}

	// 3. Save to Cache
	if redisClient != nil {
		data, _ := json.Marshal(resp)
		redisClient.Set(context.Background(), cacheKey, data, time.Duration(cfg.Interpolation.CacheTTL)*time.Second)
	}

	// 4. Publish Event
	go publishEvent(req, resY)

	c.JSON(http.StatusOK, resp)
}

func main() {
	setupRedis()
	go setupRabbitMQ()
	defer func() {
		if amqpConn != nil {
			amqpConn.Close()
		}
	}()

	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		v1.POST("/interpolate", handleInterpolate)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Interpolation.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()

	log.Info().Msgf("Interpolation Service started on port %d", cfg.Interpolation.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}

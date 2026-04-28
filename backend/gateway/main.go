package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/MindlessMuse666/interpolation/backend/gateway/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	mredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// @title Interpolation API
// @version 1.0
// @description API Gateway для учебного проекта по интерполяции.
// @contact.name API Support
// @contact.url http://localhost/support
// @contact.email mindlessmuse.666@gmail.com
// @host localhost:8080
// @BasePath /api/v1

type Config struct {
	Gateway struct {
		Port             int    `toml:"port"`
		InterpolationURL string `toml:"interpolation_url"`
		HistoryURL       string `toml:"history_url"`
		RedisAddr        string `toml:"redis_addr"`
	} `toml:"gateway"`
}

var cfg Config

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if _, err := toml.DecodeFile("../../config/config.toml", &cfg); err != nil {
		if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
			log.Warn().Err(err).Msg("Could not load config.toml, using defaults")
			cfg.Gateway.Port = 8080
			cfg.Gateway.InterpolationURL = "http://localhost:8081"
			cfg.Gateway.HistoryURL = "http://localhost:8082"
			cfg.Gateway.RedisAddr = "localhost:6379"
		}
	}
}

func proxyHandler(target string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: targetURL.Scheme,
		Host:   targetURL.Host,
	})

	// Add transport with reasonable timeouts to avoid status 502 on slow connections
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return func(c *gin.Context) {
		c.Request.URL.Path = targetURL.Path
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// @Summary Список доступных методов
// @Description Возвращает список поддерживаемых методов интерполяции (linear, lagrange, newton).
// @Produce json
// @Success 200 {array} string
// @Router /methods [get]
func handleMethods(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"linear", "lagrange", "newton"})
}

func setupRateLimiter() gin.HandlerFunc {
	option, err := redis.ParseURL(fmt.Sprintf("redis://%s", cfg.Gateway.RedisAddr))
	var client *redis.Client
	if err != nil {
		log.Warn().Msg("Redis URL parse failed, using default addr")
		client = redis.NewClient(&redis.Options{
			Addr: cfg.Gateway.RedisAddr,
		})
	} else {
		client = redis.NewClient(option)
	}

	store, err := mredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix: "rate_limiter:",
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create limiter store")
		return func(c *gin.Context) { c.Next() }
	}

	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  200,
	}

	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		key := c.ClientIP()
		context, err := instance.Get(c, key)
		if err != nil {
			log.Error().Err(err).Msg("Limiter error")
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", context.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", context.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", context.Reset))

		if context.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}

// @Summary Интерполяция данных
// @Description Проксирует запрос к сервису интерполяции для вычисления значения.
// @Accept json
// @Produce json
// @Param request body InterpolateRequest true "Параметры интерполяции"
// @Success 200 {object} InterpolateResponse "Успешное вычисление"
// @Failure 400 {object} map[string]string "Неверный формат данных"
// @Router /interpolate [post]
func _swaggerInterpolateDoc() {}

// @Summary Получить историю
// @Description Проксирует запрос к сервису истории.
// @Produce json
// @Success 200 {array} CalculationRecord "Список записей истории"
// @Router /history [get]
func _swaggerHistoryGetDoc() {}

// @Summary Очистить историю
// @Description Удаляет все записи из базы данных истории.
// @Produce json
// @Success 200 {object} map[string]string "История успешно очищена"
// @Router /history [delete]
func _swaggerHistoryDeleteDoc() {}

type InterpolateRequest struct {
	Method  string  `json:"method" example:"linear"`
	Points  []Point `json:"points"`
	TargetX float64 `json:"target_x" example:"1.5"`
}

type Point struct {
	X float64 `json:"x" example:"0.0"`
	Y float64 `json:"y" example:"0.0"`
}

type InterpolateResponse struct {
	Method string  `json:"method" example:"linear"`
	Result float64 `json:"result" example:"1.5"`
	Curve  []Point `json:"curve"`
	Cached bool    `json:"cached" example:"false"`
}

type CalculationRecord struct {
	ID        int       `json:"id" example:"1"`
	Method    string    `json:"method" example:"linear"`
	Points    []Point   `json:"points"`
	TargetX   float64   `json:"target_x" example:"1.5"`
	Result    float64   `json:"result" example:"1.5"`
	CreatedAt time.Time `json:"created_at" example:"2026-04-28T18:32:04Z"`
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	// Logger middleware
	r.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		log.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", c.Writer.Status()).
			Dur("duration", time.Since(start)).
			Msg("Request")
	})

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust for production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate Limiting
	r.Use(setupRateLimiter())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"service": "Interpolation API Gateway",
				"status":  "running",
				"version": "1.0",
			})
		})
		v1.GET("/methods", handleMethods)
		// @Router /interpolate [post]
		v1.POST("/interpolate", proxyHandler(cfg.Gateway.InterpolationURL+"/api/v1/interpolate"))
		// @Router /history [get]
		v1.GET("/history", proxyHandler(cfg.Gateway.HistoryURL+"/api/v1/history"))
		// @Router /history [delete]
		v1.DELETE("/history", proxyHandler(cfg.Gateway.HistoryURL+"/api/v1/history"))
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Gateway.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()

	log.Info().Msgf("API Gateway started on port %d", cfg.Gateway.Port)

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

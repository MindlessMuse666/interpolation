package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ulule/limiter/v3"
	mredis "github.com/ulule/limiter/v3/drivers/store/redis"
	_ "github.com/user/interpolation/backend/gateway/docs"
)

// @title Interpolation API
// @version 1.0
// @description API Gateway for Interpolation educational project.
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
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// @Summary Get list of interpolation methods
// @Description Returns supported interpolation methods
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
		v1.GET("/methods", handleMethods)
		v1.POST("/interpolate", proxyHandler(cfg.Gateway.InterpolationURL))
		v1.GET("/history", proxyHandler(cfg.Gateway.HistoryURL))
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

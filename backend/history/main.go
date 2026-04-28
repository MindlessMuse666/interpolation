package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/MindlessMuse666/interpolation/backend/core/interpolation"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type Config struct {
	History struct {
		Port         int    `toml:"port"`
		SQLitePath   string `toml:"sqlite_path"`
		RabbitMQURL  string `toml:"rabbitmq_url"`
		QueueName    string `toml:"queue_name"`
		ExchangeName string `toml:"exchange_name"`
		RoutingKey   string `toml:"routing_key"`
	} `toml:"history"`
}

type CalculationRecord struct {
	ID        int                   `json:"id"`
	Method    string                `json:"method"`
	Points    []interpolation.Point `json:"points"`
	TargetX   float64               `json:"target_x"`
	Result    float64               `json:"result"`
	CreatedAt time.Time             `json:"created_at"`
}

type CalculationEvent struct {
	Method    string                `json:"method"`
	Points    []interpolation.Point `json:"points"`
	TargetX   float64               `json:"target_x"`
	Result    float64               `json:"result"`
	Timestamp time.Time             `json:"timestamp"`
}

var (
	cfg      Config
	db       *sql.DB
	amqpConn *amqp.Connection
	amqpChan *amqp.Channel
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if _, err := toml.DecodeFile("../../config/config.toml", &cfg); err != nil {
		if _, err := toml.DecodeFile("config.toml", &cfg); err != nil {
			log.Warn().Err(err).Msg("Could not load config.toml, using defaults")
			cfg.History.Port = 8082
			cfg.History.SQLitePath = "./data/history.db"
			cfg.History.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
			cfg.History.QueueName = "history_queue"
			cfg.History.ExchangeName = "calculations"
			cfg.History.RoutingKey = "calculation.completed"
		}
	}
}

func setupDB() {
	// Create directory if not exists
	dir := filepath.Dir(cfg.History.SQLitePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal().Err(err).Msg("Failed to create database directory")
	}

	var err error
	db, err = sql.Open("sqlite3", cfg.History.SQLitePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open SQLite database")
	}

	query := `
	CREATE TABLE IF NOT EXISTS calculations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		method TEXT,
		points TEXT,
		target_x REAL,
		result REAL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatal().Err(err).Msg("Failed to create table")
	}
	log.Info().Msg("Database setup complete")
}

func setupRabbitMQ() {
	var err error
	// Retry connection to RabbitMQ
	for i := 0; i < 10; i++ {
		amqpConn, err = amqp.Dial(cfg.History.RabbitMQURL)
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
		cfg.History.ExchangeName, // name
		"direct",                 // type
		true,                     // durable
		false,                    // auto-deleted
		false,                    // internal
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to declare exchange")
	}

	q, err := amqpChan.QueueDeclare(
		cfg.History.QueueName, // name
		true,                  // durable
		false,                 // delete when unused
		false,                 // exclusive
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to declare queue")
	}

	err = amqpChan.QueueBind(
		q.Name,                   // queue name
		cfg.History.RoutingKey,   // routing key
		cfg.History.ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to bind queue")
	}

	msgs, err := amqpChan.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (we use manual ack)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register consumer")
	}

	go func() {
		for d := range msgs {
			var event CalculationEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Error().Err(err).Msg("Failed to unmarshal event")
				d.Nack(false, false)
				continue
			}

			pointsJSON, _ := json.Marshal(event.Points)
			_, err := db.Exec(
				"INSERT INTO calculations (method, points, target_x, result, created_at) VALUES (?, ?, ?, ?, ?)",
				event.Method, string(pointsJSON), event.TargetX, event.Result, event.Timestamp,
			)

			if err != nil {
				log.Error().Err(err).Msg("Failed to save calculation to DB")
				d.Nack(false, true) // requeue if DB is down?
			} else {
				d.Ack(false)
				log.Info().Str("method", event.Method).Msg("Calculation saved to history")
			}
		}
	}()

	log.Info().Msg("RabbitMQ consumer started")
}

// @Summary Получить историю вычислений
// @Description Возвращает последние 10 записей из истории интерполяций.
// @Produce json
// @Success 200 {array} CalculationRecord
// @Failure 500 {object} map[string]string "Ошибка при получении истории"
// @Router /history [get]
func handleGetHistory(c *gin.Context) {
	rows, err := db.Query("SELECT id, method, points, target_x, result, created_at FROM calculations ORDER BY created_at DESC LIMIT 10")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}
	defer rows.Close()

	var history []CalculationRecord
	for rows.Next() {
		var r CalculationRecord
		var pointsStr string
		if err := rows.Scan(&r.ID, &r.Method, &pointsStr, &r.TargetX, &r.Result, &r.CreatedAt); err != nil {
			continue
		}
		json.Unmarshal([]byte(pointsStr), &r.Points)
		history = append(history, r)
	}

	// Return empty array instead of null
	if history == nil {
		history = []CalculationRecord{}
	}

	c.JSON(http.StatusOK, history)
}

// @Summary Очистить историю
// @Description Удаляет все записи из таблицы calculations.
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /history [delete]
func handleClearHistory(c *gin.Context) {
	_, err := db.Exec("DELETE FROM calculations")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear history"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "History cleared"})
}

func main() {
	setupDB()
	go setupRabbitMQ()
	defer func() {
		if amqpConn != nil {
			amqpConn.Close()
		}
		if db != nil {
			db.Close()
		}
	}()

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/history", handleGetHistory)
		v1.DELETE("/history", handleClearHistory)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.History.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("listen")
		}
	}()

	log.Info().Msgf("History Service started on port %d", cfg.History.Port)

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

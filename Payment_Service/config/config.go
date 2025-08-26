package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App    App
		Db     Db
		Redis  Redis
		JWT    JWT
		Kafka  Kafka
		Outbox Outbox
	}

	App struct {
		Port string
	}

	Db struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}

	Redis struct {
		Addr string
	}

	Kafka struct {
		Brokers string
	}

	Outbox struct {
		Batch    int
		Interval string
		MaxRetry int
	}

	JWT struct {
		SecretKey            string
		AccessTokenDuration  int64
		RefreshTokenDuration int64
	}
)

func NewConfig() *Config {

	if err := godotenv.Load("./env/.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &Config{
		App: App{
			Port: os.Getenv("APP_PORT"),
		},
		Db: Db{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
		},
		Redis: Redis{
			Addr: os.Getenv("REDIS_ADDR"),
		},
		Kafka: Kafka{
			Brokers: os.Getenv("KAFKA_BROKERS"),
		},
		Outbox: Outbox{
			Batch: func() int {
				duration, err := strconv.ParseInt(os.Getenv("OUTBOX_BATCH"), 10, 64)
				if err != nil {
					log.Fatalf("Error getting OUTBOX_BATCH: %v", err)
				}
				return int(duration)

			}(),
			Interval: os.Getenv("OUTBOX_INTERVAL"),
			MaxRetry: func() int {
				duration, err := strconv.ParseInt(os.Getenv("OUTBOX_MAX_RETRY"), 10, 64)
				if err != nil {
					log.Fatalf("Error getting OUTBOX_MAX_RETRY: %v", err)
				}
				return int(duration)

			}(),
		},
		JWT: JWT{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
			AccessTokenDuration: func() int64 {
				duration, err := strconv.ParseInt(os.Getenv("ACCESS_TOKEN_DURATION"), 10, 64)
				if err != nil {
					log.Fatalf("Error getting ACCESS_TOKEN_DURATION: %v", err)
				}
				return int64(duration)

			}(),
			RefreshTokenDuration: func() int64 {
				duration, err := strconv.ParseInt(os.Getenv("REFRESH_TOKEN_DURATION"), 10, 64)
				if err != nil {
					log.Fatalf("Error getting REFRESH_TOKEN_DURATION: %v", err)
				}
				return int64(duration)

			}(),
		},
	}
}

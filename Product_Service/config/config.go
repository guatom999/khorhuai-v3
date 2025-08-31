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
		JWT    JWT
		Expire Expire
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
	Expire struct {
		Interval int64
		Batch    int64
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
		Expire: Expire{
			Interval: func() int64 {
				interval, err := strconv.ParseInt(os.Getenv("EXPIRE_SWEEP_INTERVAL"), 10, 64)
				if err != nil {
					log.Fatalf("Error getting EXPIRE_SWEEP_INTERVAL: %v", err)
				}
				return interval

			}(),
			Batch: func() int64 {
				batch, err := strconv.ParseInt(os.Getenv("EXPIRE_SWEEP_BATCH"), 10, 64)
				if err != nil {
					log.Fatalf("Error getting EXPIRE_SWEEP_BATCH: %v", err)
				}
				return batch

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

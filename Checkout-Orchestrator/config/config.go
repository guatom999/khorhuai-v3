package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		AppPort          string
		TemporalHostPort string
		TaskQueue        string
		InventoryBaseURL string
		OrderBaseURL     string
		PaymentBaseURL   string
	}
)

func NewConfig() *Config {

	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = "../../env/.env"
	}
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s", envPath)
	}

	return &Config{
		AppPort:          os.Getenv("APP_PORT"),
		TemporalHostPort: os.Getenv("TEMPORAL_HOSTPORT"),
		TaskQueue:        os.Getenv("CHECKOUT_TASK_QUEUE"),
		InventoryBaseURL: os.Getenv("INVENTORY_BASE_URL"),
		OrderBaseURL:     os.Getenv("ORDER_BASE_URL"),
		PaymentBaseURL:   os.Getenv("PAYMENT_BASE_URL"),
	}
}

package config

import "os"

type (
	Config struct {
		TemporalHostPort string
		TaskQueue        string
		InventoryBaseURL string
		OrderBaseURL     string
		PaymentBaseURL   string
	}
)

func NewConfig() *Config {
	return &Config{
		TemporalHostPort: os.Getenv("TEMPORAL_HOSTPORT"),
		TaskQueue:        os.Getenv("CHECKOUT_TASK_QUEUE"),
		InventoryBaseURL: os.Getenv("INVENTORY_BASE_URL"),
		OrderBaseURL:     os.Getenv("ORDER_BASE_URL"),
		PaymentBaseURL:   os.Getenv("PAYMENT_BASE_URL"),
	}
}

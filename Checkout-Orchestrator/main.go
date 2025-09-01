package main

import (
	"log"

	"github.com/guatom999/ecommerce-orchestrator/config"
	"go.temporal.io/sdk/client"
)

func main() {

	cfg := config.NewConfig()

	client, err := client.Dial(client.Options{
		HostPort: cfg.TemporalHostPort,
	})
	if err != nil {
		log.Fatalf("temporal client failed: %v", err)
	}
	defer client.Close()

}

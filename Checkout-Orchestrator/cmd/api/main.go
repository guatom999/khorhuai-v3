package main

import (
	"context"
	"log"

	"github.com/guatom999/ecommerce-orchestrator/config"
	"github.com/guatom999/ecommerce-orchestrator/server"
	"go.temporal.io/sdk/client"
)

func main() {

	ctx := context.Background()

	cfg := config.NewConfig()

	client, err := client.Dial(client.Options{
		HostPort: cfg.TemporalHostPort,
	})
	if err != nil {
		log.Fatalf("temporal client failed: %v", err)
	}
	defer client.Close()

	server.NewServer(cfg, client).Start(ctx)

}

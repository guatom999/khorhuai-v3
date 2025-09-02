package main

import (
	"log"

	"github.com/guatom999/ecommerce-orchestrator/activities"
	"github.com/guatom999/ecommerce-orchestrator/config"
	"github.com/guatom999/ecommerce-orchestrator/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
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

	acts := activities.NewActitivities(cfg)

	w := worker.New(client, cfg.TaskQueue, worker.Options{})

	w.RegisterWorkflow(workflows.CheckoutWorkflow)

	w.RegisterActivity(acts.ReserveStock)
	w.RegisterActivity(acts.ReleaseStock)
	w.RegisterActivity(acts.CommitStock)
	w.RegisterActivity(acts.CreateOrder)
	w.RegisterActivity(acts.CancelOrder)
	w.RegisterActivity(acts.ConfirmOrder)
	w.RegisterActivity(acts.CreatePayment)

}

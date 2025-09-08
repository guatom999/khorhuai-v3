package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guatom999/ecommerce-orchestrator/config"
	"github.com/guatom999/ecommerce-orchestrator/workflows"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.temporal.io/sdk/client"
)

type (
	server struct {
		app    *echo.Echo
		cfg    *config.Config
		client client.Client
	}

	Item struct {
		ProductID string `json:"product_id"`
		Title     string `json:"title"`
		UnitPrice int    `json:"unit_price"`
		Quantity  int    `json:"quantity"`
	}

	Request struct {
		OrderID     string `json:"order_id"`
		UserID      string `json:"user_id"`
		Currency    string `json:"currency"`
		AmountCents int64  `json:"amount_cents"`
		TTLSeconds  int    `json:"ttl_seconds"`
		Items       []Item `json:"items"`
	}
)

func NewServer(cfg *config.Config, client client.Client) *server {
	return &server{
		app:    echo.New(),
		cfg:    cfg,
		client: client,
	}
}

func (s *server) gratefulShutdown(pctx context.Context, close <-chan os.Signal) {
	<-close

	ctx, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to Shutting down server: %v", err)
	}

	log.Println("Shutting Down Server...")
}

func (s *server) Start(pctx context.Context) {

	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Error: Request Timeout",
		Timeout:      time.Second * 10,
	}))

	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
	}))

	s.app.Use(middleware.Logger())

	s.app.POST("/checkout", func(c echo.Context) error {

		req := new(Request)

		if err := c.Bind(req); err != nil || len(req.Items) == 0 || req.OrderID == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid body"})
		}

		in := workflows.CheckoutInput{
			OrderID:     req.OrderID,
			UserID:      req.UserID,
			Currency:    req.Currency,
			AmountCents: req.AmountCents,
			TTLSeconds:  req.TTLSeconds,
		}

		for _, item := range req.Items {
			in.Items = append(in.Items, workflows.Item{
				ProductID: item.ProductID,
				Title:     item.Title,
				UnitPrice: item.UnitPrice,
				Quantity:  item.Quantity,
			})
		}

		run, err := s.client.ExecuteWorkflow(c.Request().Context(),
			client.StartWorkflowOptions{
				ID:        "checkout_" + req.OrderID,
				TaskQueue: s.cfg.TaskQueue,
			},
			workflows.CheckoutWorkflow, in,
		)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusAccepted, echo.Map{
			"workflow_id": "checkout_" + req.OrderID,
			"run_id":      run.GetRunID()})

	})

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)
	go s.gratefulShutdown(pctx, close)

	fmt.Println("app port is ", s.cfg.AppPort)

	if err := s.app.Start(s.cfg.AppPort); err != nil {
		log.Fatalf("Failed to shutdown:%v", err)
	}

}

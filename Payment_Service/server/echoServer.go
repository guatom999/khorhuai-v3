package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/guatom999/ecommerce-payment-api/config"
	paymenthandlers "github.com/guatom999/ecommerce-payment-api/modules/paymentHandlers"
	paymentrepositories "github.com/guatom999/ecommerce-payment-api/modules/paymentRepositories"
	paymentusecases "github.com/guatom999/ecommerce-payment-api/modules/paymentUseCases"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	server struct {
		app *echo.Echo
		db  *sqlx.DB
		cfg *config.Config
	}
)

func NewEchoServer(cfg *config.Config, db *sqlx.DB) *server {
	return &server{
		app: echo.New(),
		cfg: cfg,
		db:  db,
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

	s.paymentModules()

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)
	go s.gratefulShutdown(pctx, close)

	if err := s.app.Start(s.cfg.App.Port); err != nil {
		log.Fatalf("Failed to shutdown:%v", err)
	}

}

func (s *server) paymentModules() {
	paymentRepo := paymentrepositories.NewPaymentRepository(s.db)
	paymentUsecase := paymentusecases.NewPaymentUsecase(paymentRepo)
	paymenthandlers := paymenthandlers.NewPaymenthandler(paymentUsecase)

	paymentRoute := s.app.Group("/app/v1/payments")

	paymentRoute.POST("/", paymenthandlers.CreatePayment)

	paymentRoute.GET("/:id", paymenthandlers.GetPayment)

	paymentRoute.PATCH("/:id", paymenthandlers.UpdatePaymentStatus)

}

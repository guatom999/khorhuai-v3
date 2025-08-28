package main

import (
	"context"

	"github.com/guatom999/ecommerce-payment-api/config"
	"github.com/guatom999/ecommerce-payment-api/databases"
	"github.com/guatom999/ecommerce-payment-api/databases/redisdb"
	"github.com/guatom999/ecommerce-payment-api/server"
	"github.com/guatom999/ecommerce-payment-api/utils"
	"go.uber.org/zap"
)

func main() {

	ctx := context.Background()

	cfg := config.NewConfig()

	db := databases.ConnDB(cfg)

	redisDb := redisdb.NewRedis(cfg)

	utils.InitLogger()
	defer utils.ShutdownLogger()

	shutdown, err := utils.InitTracing(ctx, utils.OtelConfig{
		ServiceName: "payment-api",
		Endpoint:    cfg.Otel.Endpoint,
		SampleRatio: 1.0,
	})
	if err != nil {
		utils.AppLogger().Fatal("init tracing", zap.Error(err))
	}
	defer shutdown(context.Background())

	server.NewEchoServer(cfg, db, redisDb).Start(ctx)

	utils.AppLogger().Info("payment api starting", zap.String("addr", ":8082"))
}

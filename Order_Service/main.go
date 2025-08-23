package main

import (
	"context"

	"github.com/guatom999/ecommerce-order-api/config"
	"github.com/guatom999/ecommerce-order-api/databases"
	"github.com/guatom999/ecommerce-order-api/server"
)

func main() {

	ctx := context.Background()

	cfg := config.NewConfig()

	db := databases.ConnDB(cfg)

	server.NewEchoServer(cfg, db).Start(ctx)

}

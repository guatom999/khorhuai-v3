package main

import (
	"context"

	"github.com/guatom999/ecommerce-notification-api/config"
	"github.com/guatom999/ecommerce-notification-api/databases"
	"github.com/guatom999/ecommerce-notification-api/server"
)

func main() {

	ctx := context.Background()

	cfg := config.NewConfig()

	db := databases.ConnDB(cfg)

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr: cfg.Redis.Addr,
	// })

	// redisDb := redisdb.Store{Rdb: rdb}
	// _ = redisDb

	server.NewEchoServer(cfg, db).Start(ctx)
}

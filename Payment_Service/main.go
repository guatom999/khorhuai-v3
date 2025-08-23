package main

import (
	"context"

	"github.com/guatom999/ecommerce-payment-api/config"
	"github.com/guatom999/ecommerce-payment-api/databases"
	redisdb "github.com/guatom999/ecommerce-payment-api/databases/redisdb"
	"github.com/guatom999/ecommerce-payment-api/server"
	"github.com/redis/go-redis/v9"
)

func main() {

	ctx := context.Background()

	cfg := config.NewConfig()

	db := databases.ConnDB(cfg)

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.Redis.Addr,
	})

	redisDb := redisdb.Store{Rdb: rdb}
	_ = redisDb

	server.NewEchoServer(cfg, db).Start(ctx)

}

package databases

import (
	"context"
	"fmt"
	"log"

	"github.com/XSAM/otelsql"
	"github.com/guatom999/ecommerce-payment-api/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnDB(cfg *config.Config) *sqlx.DB {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Password, cfg.Db.DBName,
	)

	driverName, err := otelsql.Register("postgres", otelsql.WithAttributes(), otelsql.WithSpanNameFormatter(func(ctx context.Context, method otelsql.Method, query string) string {
		return fmt.Sprintf("db.%s", method)
	}))

	if err != nil {
		log.Printf("Error registering otelsql driver: %v", err)
		panic(err)
	}

	// db, err := sqlx.Connect("postgres", connStr)
	db, err := sqlx.Open(driverName, connStr)
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error connecting to the database: %v", err)
		panic(err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db

}

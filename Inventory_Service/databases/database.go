package databases

import (
	"fmt"
	"log"

	"github.com/guatom999/ecommerce-product-api/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func ConnDB(cfg *config.Config) *sqlx.DB {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Db.Host, cfg.Db.Port, cfg.Db.User, cfg.Db.Password, cfg.Db.DBName,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Printf("Error connecting to the database: %v", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		log.Printf("Error connecting to the database: %v", err)
		panic(err)
	}

	return db

}

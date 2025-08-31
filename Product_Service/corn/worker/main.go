package main

import (
	"context"
	"log"
	"time"

	"github.com/guatom999/ecommerce-product-api/config"
	"github.com/guatom999/ecommerce-product-api/databases"
	"github.com/jmoiron/sqlx"
)

type (
	item struct {
		ProductID string `db:"product_id"`
		Quantity  int    `db:"quantity"`
	}
)

func main() {
	ctx := context.Background()

	cfg := config.NewConfig()

	db := databases.ConnDB(cfg)
	defer db.Close()

	log.Printf("stock-expirer started: interval=%v batch=%d", cfg.Expire.Interval, cfg.Expire.Batch)

	ticker := time.NewTicker(time.Duration(cfg.Expire.Interval))
	defer ticker.Stop()

	for {
		if err := sweepOnce(ctx, db, int(cfg.Expire.Batch)); err != nil {
			log.Printf("sweep error: %v", err)
		}
		<-ticker.C
	}

}

func sweepOnce(ctx context.Context, db *sqlx.DB, batch int) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	ids := make([]string, 0)

	if err := tx.SelectContext(ctx, &ids,
		`
		SELECT id
		FROM stock_reservations
		WHERE status = 'held'
		  AND expires_at < CURRENT_TIMESTAMP
		FOR UPDATE SKIP LOCKED
		LIMIT $1
		`, batch,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, id := range ids {
		if err := releasereservation(ctx, db, id); err != nil {
			log.Printf("release reservation failed %v %v", id, err)
		}
	}

	return nil

}

func releasereservation(ctx context.Context, db *sqlx.DB, reserveId string) error {

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var status string
	if err := tx.GetContext(ctx, &status, `SELECT status FROM stock_reservations WHERE id = $1 FOR UPDATE`, reserveId); err != nil {
		return err
	}
	if status != "held" {
		return tx.Commit()
	}

	items := make([]item, 0)

	if err := tx.SelectContext(ctx, &items, `SELECT product_id, quantity FROM stock_reservation_items WHERE reservation_id = $1`, reserveId); err != nil {
		return err
	}

	for _, item := range items {
		if _, err := tx.ExecContext(ctx, `
					UPDATE products
			SET stock_qty = stock_qty + $2,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
			`, item.ProductID, item.Quantity,
		); err != nil {
			return err
		}
	}

	if _, err = tx.ExecContext(ctx, `
		UPDATE stock_reservations
		SET status='released', updated_at=CURRENT_TIMESTAMP
		WHERE id=$1 AND status='held'
	`, reserveId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}

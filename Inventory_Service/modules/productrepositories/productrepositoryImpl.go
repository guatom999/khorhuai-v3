package productrepositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/guatom999/ecommerce-product-api/modules"
	"github.com/jmoiron/sqlx"
)

type (
	productRepositoryImpl struct {
		db *sqlx.DB
	}
)

func NewPortRepository(db *sqlx.DB) ProductRepositoryInterface {
	return &productRepositoryImpl{db: db}
}

func (r *productRepositoryImpl) Reserve(ctx context.Context, input modules.ReserveInput) (string, error) {

	if len(input.Items) == 0 {
		return "", errors.New("items required")
	}
	if input.TTL <= 0 {
		input.TTL = 15 * time.Minute
	}

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// create reservation header
	var reserveID string
	err = tx.GetContext(ctx, &reserveID, `
	  INSERT INTO stock_reservations (order_id, user_id, expires_at)
	  VALUES (NULLIF($1,'' )::uuid, NULLIF($2,'' )::uuid, CURRENT_TIMESTAMP + make_interval(secs => $3))
	  RETURNING id;
	`, input.OrderID, input.UserID, int(input.TTL.Seconds()))
	if err != nil {
		return "", err
	}

	for _, it := range input.Items {
		res, err := tx.ExecContext(ctx, `
		  UPDATE products
		  SET stock_qty = stock_qty - $2,
		      updated_at = CURRENT_TIMESTAMP
		  WHERE id = $1 AND stock_qty >= $2
		`, it.ProductID, it.Quantity)
		if err != nil {
			return "", err
		}
		aff, _ := res.RowsAffected()
		if aff == 0 {
			err = errors.New("insufficient stock for product " + it.ProductID)
			return "", err
		}

		// 2) append to reservation_items
		if _, err = tx.ExecContext(ctx, `
		  INSERT INTO stock_reservation_items (reservation_id, product_id, quantity)
		  VALUES ($1, $2, $3)
		`, reserveID, it.ProductID, it.Quantity); err != nil {
			return "", err
		}
	}

	err = tx.Commit()

	return reserveID, err
}
func (r *productRepositoryImpl) Release(ctx context.Context, reservationID string) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var status string
	if err := r.db.GetContext(ctx, &status, `SELECT status FROM stock_reservations WHERE id = $1`, reservationID); err != nil {
		return err
	}

	if status == "released" {
		return nil
	}
	if status == "committed" {
		return nil
	}

	items := make([]modules.Item, 0)
	if err := tx.SelectContext(ctx, &items, `SELECT product_id, quantity FROM stock_reservation_items WHERE reservation_id=$1`, reservationID); err != nil {
		return err
	}

	for _, v := range items {
		if _, err := tx.ExecContext(ctx, `
			UPDATE products
		  	SET stock_qty = stock_qty + $2,
		    	updated_at = CURRENT_TIMESTAMP
		  	WHERE id = $1
		`, v.ProductId, v.Quantity); err != nil {
			return err
		}
	}

	if _, err = tx.ExecContext(ctx, `
	  UPDATE stock_reservations
	  SET status='released', updated_at=CURRENT_TIMESTAMP
	  WHERE id=$1 AND status='held'
	`, reservationID); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
func (r *productRepositoryImpl) Commit(ctx context.Context, reservationID string) error {
	if _, err := r.db.ExecContext(ctx, `
	  UPDATE stock_reservations
	  SET status='committed', updated_at=CURRENT_TIMESTAMP
	  WHERE id=$1 AND status='held'
	`, reservationID); err != nil {
		return err
	}

	return nil
}

package paymentrepositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/guatom999/ecommerce-payment-api/modules"
	"github.com/jmoiron/sqlx"
)

type (
	paymentRepository struct {
		db *sqlx.DB
	}
)

func NewPaymentRepository(db *sqlx.DB) PaymentRepositoryInterface {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreateProcessing(ctx context.Context, in modules.CreatePaymentRequest) (string, error) {

	var userPtr *string

	if in.UserID != "" {
		userPtr = &in.UserID
	}

	var id string
	err := r.db.GetContext(ctx, &id, `
	  INSERT INTO payments (order_id, user_id, amount, currency, status)
	  VALUES ($1, $2::uuid, $3, $4, 'processing')
	  RETURNING id;
	`, in.OrderID, userPtr, in.Amount, in.Currency)
	return id, err

}
func (r *paymentRepository) UpdateStatus(ctx context.Context, id, status string) error {

	_, err := r.db.ExecContext(ctx, `
	  UPDATE payments
	  SET status=$2, updated_at=CURRENT_TIMESTAMP
	  WHERE id=$1
	`, id, status)

	if err != nil {
		return err
	}

	return nil
}
func (r *paymentRepository) Get(ctx context.Context, id string) (*modules.PaymentRow, error) {

	p := new(modules.PaymentRow)

	err := r.db.GetContext(ctx, p, `
	  SELECT
	    id, order_id,
	    COALESCE(user_id::text,'') AS user_id_text,
	    amount, currency, status, created_at, updated_at
	  FROM payments
	  WHERE id = $1
	`, id)

	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *paymentRepository) UpdateStatusWithOutbox(ctx context.Context, id, newStatus string) error {

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	res, err := tx.ExecContext(ctx, `
	  UPDATE payments
	  SET status=$2, updated_at=CURRENT_TIMESTAMP
	  WHERE id=$1
	`, id, newStatus,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if rows == 0 {
		return errors.New("payment not found")
	}

	payload := modules.PaymentStatusChanged{
		EventType:  "payment.status.changed",
		PaymentID:  id,
		NewStatus:  newStatus,
		OccurredAt: time.Now(),
	}
	b, _ := json.Marshal(payload)

	if _, err = r.db.ExecContext(ctx, `
		INSERT INTO outbox_events (aggregate_type, aggregate_id, topic, key, payload, status)
	    VALUES ($1, $2::uuid, $3, $4, $5, 'pending')
	`,
		"payment", id, "payment.events", id, b,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error: failed to commit sql %v", err.Error())
		return err
	}

	return nil
}

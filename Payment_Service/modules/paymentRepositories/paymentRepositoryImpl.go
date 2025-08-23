package paymentrepositories

import (
	"context"

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

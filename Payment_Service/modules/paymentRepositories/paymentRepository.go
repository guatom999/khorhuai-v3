package paymentrepositories

import (
	"context"

	"github.com/guatom999/ecommerce-payment-api/modules"
)

type PaymentRepositoryInterface interface {
	CreateProcessing(ctx context.Context, in modules.CreatePaymentInput) (string, error)
	UpdateStatus(ctx context.Context, id, status string) error
	Get(ctx context.Context, id string) (*modules.PaymentRow, error)
}

package paymentusecases

import (
	"context"

	"github.com/guatom999/ecommerce-payment-api/modules"
)

type PaymentUsecaseInterface interface {
	CreatePayment(ctx context.Context, cmd modules.CreatePaymentCommand) (*modules.PaymentRow, error)
	GetPayment(ctx context.Context, id string) (*modules.PaymentRow, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

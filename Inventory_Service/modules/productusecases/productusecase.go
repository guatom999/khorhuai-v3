package productusecases

import (
	"context"

	"github.com/guatom999/ecommerce-product-api/modules"
)

type (
	ProductusecaseInterface interface {
		Reserve(ctx context.Context, input modules.ReserveInput) (string, error)
		Release(ctx context.Context, reservationID string) error
		Commit(ctx context.Context, reservationID string) error
	}
)

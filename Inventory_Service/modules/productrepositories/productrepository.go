package productrepositories

import (
	"context"

	"github.com/guatom999/ecommerce-product-api/modules"
)

type ProductRepositoryInterface interface {
	Reserve(ctx context.Context, input modules.ReserveInput) (string, error)
	ReleaseExpired(ctx context.Context) error
	Release(ctx context.Context, reservationID string) error
	Commit(ctx context.Context, reservationID string) error
}

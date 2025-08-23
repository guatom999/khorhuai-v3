package cartUsecases

import (
	"context"

	"github.com/guatom999/ecommerce-shopping-cart-api/modules"
)

type (
	CartUsecaseInterface interface {
		GetCart(ctx context.Context, userID string, sessionID string) (*modules.CartView, error)
		UpsertItem(ctx context.Context, userID string, sessionID string, productID string, qty int, unitPrice int64, currency string) (*modules.CartView, error)
		Merge(ctx context.Context, userID string, sessionID string) error
	}
)

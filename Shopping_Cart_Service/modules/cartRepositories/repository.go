package cartRepositories

import (
	"context"

	"github.com/guatom999/ecommerce-shopping-cart-api/modules"
)

type (
	CartRepositoryInterface interface {
		ResolveCartID(ctx context.Context, userID string, sessionID string) (string, error)
		UpsertCartItem(ctx context.Context, cartID, productID string, qty int, unitPrice int64, currency string) error
		MergeGuestToUserCart(ctx context.Context, userID, sessionID string) error
		GetCart(ctx context.Context, userID string, sessionID string) (*modules.CartView, error)
	}
)

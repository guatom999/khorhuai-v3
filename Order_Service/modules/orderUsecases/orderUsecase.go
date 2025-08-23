package orderUsecases

import (
	"context"

	"github.com/guatom999/ecommerce-order-api/modules"
)

type (
	OrderUsecaseInterface interface {
		Create(ctx context.Context, in modules.CreateOrderInput) (string, error)
		Get(ctx context.Context, orderID string) (*modules.OrderDetail, error)
		ListByUser(ctx context.Context, userID string, limit, offset int) ([]modules.OrderListItem, error)
		UpdateStatus(ctx context.Context, orderID, status string) error
	}
)

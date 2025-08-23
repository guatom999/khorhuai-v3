package orderRepositories

import (
	"context"

	"github.com/guatom999/ecommerce-order-api/modules"
)

type (
	OrderRepositoryInterface interface {
		CreateOrder(pctx context.Context, in modules.CreateOrderInput) (string, error)
		GetOrder(pctx context.Context, orderID string) (*modules.OrderDetail, error)
		ListOrdersByUser(pctx context.Context, userID string, limit, offset int) ([]modules.OrderListItem, error)
		UpdateOrderStatus(pctx context.Context, orderID string, status string) error
	}
)

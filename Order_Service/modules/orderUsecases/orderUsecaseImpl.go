package orderUsecases

import (
	"context"

	"github.com/guatom999/ecommerce-order-api/modules"
	"github.com/guatom999/ecommerce-order-api/modules/orderRepositories"
)

type (
	orderUsecase struct {
		orderRepo orderRepositories.OrderRepositoryInterface
	}
)

func NewOrderUsecase(orderRepo orderRepositories.OrderRepositoryInterface) OrderUsecaseInterface {
	return &orderUsecase{orderRepo: orderRepo}
}

func (u *orderUsecase) Create(pctx context.Context, in modules.CreateOrderInput) (string, error) {

	return u.orderRepo.CreateOrder(pctx, in)
}
func (u *orderUsecase) Get(pctx context.Context, orderID string) (*modules.OrderDetail, error) {

	return u.orderRepo.GetOrder(pctx, orderID)
}
func (u *orderUsecase) ListByUser(pctx context.Context, userID string, limit, offset int) ([]modules.OrderListItem, error) {

	return u.orderRepo.ListOrdersByUser(pctx, userID, limit, offset)
}
func (u *orderUsecase) UpdateStatus(pctx context.Context, orderID, status string) error {
	return u.orderRepo.UpdateOrderStatus(pctx, orderID, status)
}

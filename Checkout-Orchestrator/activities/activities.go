package activities

import (
	"context"

	"github.com/guatom999/ecommerce-orchestrator/clients"
	"github.com/guatom999/ecommerce-orchestrator/config"
	"github.com/guatom999/ecommerce-orchestrator/workflows"
)

type (
	// Config struct {
	// 	InventoryBaseUrl string
	// 	OrderBaseUrl     string
	// 	PaymentBaseUrl   string
	// }

	Activities struct {
		Inventory *clients.InventoryClient
		Order     *clients.OrderClient
		Payment   *clients.PaymentClient
	}
)

func NewActitivities(cfg *config.Config) *Activities {
	return &Activities{
		Inventory: clients.NewInventoryClient(cfg.InventoryBaseURL),
		Order:     clients.NewOrderClient(cfg.OrderBaseURL),
		Payment:   clients.NewPaymentClient(cfg.PaymentBaseURL),
	}
}

func (a *Activities) ReserveStock(ctx context.Context, orderId string, userId string, items []workflows.Item, ttl int) (string, error) {

	body := make([]map[string]any, 0, len(items))
	for _, item := range items {
		body = append(body, map[string]any{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
		})
	}

	return a.Inventory.Reserve(ctx, orderId, userId, ttl, body)
}

func (a *Activities) ReleaseStock(ctx context.Context, reservationId string) error {

	return a.Inventory.Release(ctx, reservationId)

}

func (a *Activities) CommitStock(ctx context.Context, reservationId string) error {
	return a.Inventory.Commit(ctx, reservationId)
}

func (a *Activities) CreateOrder(ctx context.Context, orderID, userID string, items []workflows.Item, currency string, amount int64) error {

	body := make([]map[string]any, 0, len(items))
	for _, item := range items {
		body = append(body, map[string]any{
			"product_id": item.ProductID,
			"quantity":   item.Quantity,
		})
	}

	return a.Order.Create(ctx, orderID, userID, currency, amount, body)
}

func (a *Activities) ConfirmOrder(ctx context.Context, orderId string) error {
	return a.Order.Confirm(ctx, orderId)
}

func (a *Activities) CancelOrder(ctx context.Context, orderId string) error {
	return a.Order.Cancel(ctx, orderId)
}

func (a *Activities) CreatePayment(ctx context.Context, orderId string, userId string, amount int64, currency string) (string, error) {
	return a.Payment.Create(ctx, orderId, userId, amount, currency)
}

package clients

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type (
	OrderClient struct {
		baseUrl string
		http    *resty.Client
	}
)

func NewOrderClient(baseUrl string) *OrderClient {
	return &OrderClient{
		baseUrl: baseUrl,
		http:    resty.New(),
	}
}

func (h *OrderClient) Create(ctx context.Context, orderID, userID, currency string, total int64, items []map[string]any) error {
	_, err := h.http.R().SetContext(ctx).
		SetBody(map[string]any{
			"id":          orderID,
			"user_id":     userID,
			"currency":    currency,
			"total_price": total,
			"items":       items,
		}).
		Post(h.baseUrl + "/app/v1/orders")
	return err
}

func (h *OrderClient) Confirm(ctx context.Context, orderId string) error {
	_, err := h.http.R().SetContext(ctx).
		Post(h.baseUrl + "/app/v1/orders/" + orderId + "/confirm")
	return err
}

func (h *OrderClient) Cancel(ctx context.Context, orderID string) error {
	_, err := h.http.R().SetContext(ctx).
		Post(h.baseUrl + "/app/v1/orders/" + orderID + "/cancel")
	return err
}

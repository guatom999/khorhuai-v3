package clients

import (
	"context"

	"github.com/go-resty/resty/v2"
)

type (
	InventoryClientInterface interface {
	}

	InventoryClient struct {
		baseUrl string
		http    *resty.Client
	}
)

func NewInventoryClient(baseUrl string) *InventoryClient {
	return &InventoryClient{
		baseUrl: baseUrl,
		http:    resty.New(),
	}
}

func (h *InventoryClient) Reserve(ctx context.Context, orderId string, ttl int, items []map[string]any) (string, error) {

	type reserveOutput struct {
		ReservationId string `json:"reservation_id"`
	}

	rsvOut := new(reserveOutput)
	if _, err := h.http.R().SetContext(ctx).SetBody(
		map[string]any{
			"order_id":    orderId,
			"ttl_seconds": ttl,
			"items":       items,
		},
	).SetResult(rsvOut).Post(h.baseUrl + "/stock/reserve"); err != nil {
		return "", err
	}

	return rsvOut.ReservationId, nil
}

func (h *InventoryClient) Release(ctx context.Context, reservationId string) error {

	if _, err := h.http.R().SetContext(ctx).SetBody(
		map[string]any{
			"reservation_id": reservationId,
		},
	).Post(h.baseUrl + "/stock/release"); err != nil {
		return err
	}

	return nil
}

func (h *InventoryClient) Commit(ctx context.Context, reservationId string) error {
	_, err := h.http.R().SetContext(ctx).
		SetBody(map[string]string{"reservation_id": reservationId}).
		Post(h.baseUrl + "/stock/commit")
	return err
}

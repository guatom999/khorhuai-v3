package modules

import "time"

type ReserveItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type ReserveStockRequest struct {
	OrderID    string        `json:"order_id"`
	UserID     string        `json:"user_id"`
	TTLSeconds int           `json:"ttl_seconds"`
	Items      []ReserveItem `json:"items"`
}

type ReleaseStockRequest struct {
	ReservationID string `json:"reservation_id"`
}

type CommitStockRequest struct {
	ReservationID string `json:"reservation_id"`
}

type ReserveInput struct {
	OrderID string
	UserID  string
	TTL     time.Duration
	Items   []ReserveItem
}

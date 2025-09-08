package modules

import "time"

type StockReservation struct {
	ID        string    `db:"id"         json:"id"`
	OrderID   *string   `db:"order_id"   json:"order_id,omitempty"`
	UserID    *string   `db:"user_id"    json:"user_id,omitempty"`
	Status    string    `db:"status"     json:"status"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type StockReservationItem struct {
	ReservationID string `db:"reservation_id" json:"reservation_id"`
	ProductID     string `db:"product_id"     json:"product_id"`
	Quantity      int    `db:"quantity"       json:"quantity"`
}

type Item struct {
	ProductId string `db:"product_id"`
	Quantity  int    `db:"quantity"`
}

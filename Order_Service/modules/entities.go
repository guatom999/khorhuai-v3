package modules

import "time"

type Order struct {
	ID         string    `db:"id"`
	UserID     *string   `db:"user_id"`
	Currency   string    `db:"currency"`
	TotalPrice int64     `db:"total_price"`
	Status     string    `db:"status"`
	Shipping   []byte    `db:"shipping_address"`
	Billing    []byte    `db:"billing_address"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type CreateOrderInput struct {
	UserID       *string
	Currency     string
	Items        []OrderItemInput
	ShippingAddr map[string]any
	BillingAddr  map[string]any
}

type OrderItemInput struct {
	ProductID string
	Title     string
	SKU       *string
	Quantity  int
	UnitPrice int64
}

type OrderDetail struct {
	ID           string         `json:"id"`
	UserID       *string        `json:"user_id,omitempty"`
	Currency     string         `json:"currency"`
	TotalPrice   int64          `json:"total_price"`
	Status       string         `json:"status"`
	ShippingAddr map[string]any `json:"shipping_address,omitempty"`
	BillingAddr  map[string]any `json:"billing_address,omitempty"`
	Items        []OrderItemRow `json:"items"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type OrderItemRow struct {
	ID         string    `json:"id"`
	ProductID  string    `json:"product_id"`
	Title      string    `json:"title"`
	SKU        *string   `json:"sku,omitempty"`
	Quantity   int       `json:"quantity"`
	UnitPrice  int64     `json:"unit_price"`
	TotalPrice int64     `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderListItem struct {
	ID         string    `json:"id"`
	TotalPrice int64     `json:"total_price"`
	Status     string    `json:"status"`
	Currency   string    `json:"currency"`
	CreatedAt  time.Time `json:"created_at"`
}

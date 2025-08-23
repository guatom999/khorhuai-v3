package modules

type CartView struct {
	CartID        string        `json:"cart_id"`
	Currency      string        `json:"currency"`
	Items         []CartItemDTO `json:"items"`
	SubtotalCents int64         `json:"subtotal_cents"`
}

type CartItemDTO struct {
	ItemID         string `json:"item_id"`
	ProductID      string `json:"product_id"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int64  `json:"unit_price_cents"`
	LineTotalCents int64  `json:"line_total_cents"`
	Currency       string `json:"currency"`
}

type UpsertCartItemRequest struct {
	ProductID      string `json:"product_id" validate:"required,uuid4"`
	Quantity       int    `json:"quantity" validate:"required,min=1"`
	UnitPriceCents int64  `json:"unit_price_cents" validate:"required,min=0"`
	Currency       string `json:"currency" validate:"required,len=3"`
}

type MergeCartRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid4"`
	SessionID string `json:"session_id" validate:"required"`
}

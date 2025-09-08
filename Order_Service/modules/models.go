package modules

type OrderItemInputRequest struct {
	ProductID string  `json:"product_id" validate:"required,uuid4"`
	Title     string  `json:"title" validate:"required"`
	SKU       *string `json:"sku,omitempty"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	UnitPrice int64   `json:"unit_price" validate:"required,min=0"`
}

type CreateOrderRequest struct {
	ID         string           `json:"id,omitempty"`
	UserID     *string          `json:"user_id,omitempty"`
	Currency   string           `json:"currency" validate:"required,len=3"`
	TotalPrice int64            `json:"total_price,omitempty"`
	Items      []OrderItemInput `json:"items" validate:"required,min=1,dive"`
}

// type CreateOrderRequest struct {
// 	UserID       *string          `json:"user_id,omitempty"`
// 	Currency     string           `json:"currency" validate:"required,len=3"`
// 	Items        []OrderItemInput `json:"items" validate:"required,min=1,dive"`
// 	ShippingAddr map[string]any   `json:"shipping_address,omitempty"`
// 	BillingAddr  map[string]any   `json:"billing_address,omitempty"`
// }

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

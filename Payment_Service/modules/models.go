package modules

type CreatePaymentRequest struct {
	OrderID  string `json:"order_id" validate:"required,uuid4"`
	UserID   string `json:"user_id,omitempty"`
	Amount   int64  `json:"amount" validate:"required,min=0"`
	Currency string `json:"currency" validate:"required,len=3"`
}

type CreatePaymentCommand struct {
	OrderID        string
	UserID         string
	Amount         int64
	Currency       string
	IdempotencyKey string
}

type UpdatePaymentStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

// type CreatePaymentInput struct {
// 	OrderID  string
// 	UserID   string
// 	Amount   int64
// 	Currency string
// }

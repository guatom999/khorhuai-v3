package modules

import "time"

type CreatePaymentRequest struct {
	OrderID  string `json:"order_id" validate:"required"`
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

type PaymentStatusChanged struct {
	EventType  string    `json:"event_type"`
	PaymentID  string    `json:"payment_id"`
	NewStatus  string    `json:"new_status"`
	OccurredAt time.Time `json:"occurred_at"`
}

type WebhookEvent struct {
	OrderId   string `json:"order_id"`
	Status    string `json:"status"` // "succeeded" or "failed"
	PaymentId string `json:"payment_id"`
	EventId   string `json:"event_id"`
}

// type CreatePaymentInput struct {
// 	OrderID  string
// 	UserID   string
// 	Amount   int64
// 	Currency string
// }

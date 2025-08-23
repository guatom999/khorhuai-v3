package modules

import "time"

type PaymentRow struct {
	ID        string    `db:"id"           json:"id"`
	OrderID   string    `db:"order_id"     json:"order_id"`
	UserID    string    `db:"user_id_text" json:"user_id,omitempty"`
	Amount    int64     `db:"amount"       json:"amount"`
	Currency  string    `db:"currency"     json:"currency"`
	Status    string    `db:"status"       json:"status"`
	CreatedAt time.Time `db:"created_at"   json:"created_at"`
	UpdatedAt time.Time `db:"updated_at"   json:"updated_at"`
}

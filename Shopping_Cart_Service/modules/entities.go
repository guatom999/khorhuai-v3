package modules

import "time"

type CartRow struct {
	ID        string    `db:"id"`
	UserID    *string   `db:"user_id"`
	SessionID *string   `db:"session_id"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ItemRow struct {
	ID             string    `db:"id"`
	ProductID      string    `db:"product_id"`
	Quantity       int       `db:"quantity"`
	UnitPriceCents int64     `db:"unit_price_cents"`
	Currency       string    `db:"currency"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

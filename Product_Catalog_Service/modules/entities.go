package modules

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`
	Price       int64     `db:"price"`
	StockQty    int       `db:"stock_qty"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Category struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Slug      string    `db:"slug"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProductCategory struct {
	ProductID  uuid.UUID `db:"product_id"`
	CategoryID uuid.UUID `db:"category_id"`
}

type ProductWithCategory struct {
	ID          uuid.UUID  `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	Price       int64      `db:"price"`
	StockQty    int        `db:"stock_qty"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	Categories  []Category `db:"categories"`
}

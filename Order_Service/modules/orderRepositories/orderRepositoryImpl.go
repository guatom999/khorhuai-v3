package orderRepositories

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/guatom999/ecommerce-order-api/modules"
	"github.com/jmoiron/sqlx"
)

type (
	orderRepository struct {
		db *sqlx.DB
	}
)

func NewOrderRepository(db *sqlx.DB) OrderRepositoryInterface {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(pctx context.Context, in modules.CreateOrderInput) (string, error) {

	if len(in.Items) == 0 {
		return "", errors.New("no items")
	}

	var total int64
	for _, item := range in.Items {
		total += int64(item.Quantity) * item.UnitPrice
	}

	shipJSON, _ := json.Marshal(in.ShippingAddr)
	billJSON, _ := json.Marshal(in.BillingAddr)

	tx, err := r.db.BeginTxx(pctx, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var orderID string
	err = tx.GetContext(pctx, &orderID, `
		INSERT INTO orders (user_id, currency, total_price, status, shipping_address, billing_address)
		VALUES ($1,$2,$3,'pending',$4,$5)
		RETURNING id;
	`, in.UserID, in.Currency, total, shipJSON, billJSON)
	if err != nil {
		return "", err
	}

	stmt, err := tx.PreparexContext(pctx, `
		INSERT INTO order_items (order_id, product_id, title, sku, quantity, unit_price, total_price)
		VALUES ($1,$2,$3,$4,$5,$6,$7);
	`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	for _, it := range in.Items {
		line := int64(it.Quantity) * it.UnitPrice
		if _, err = stmt.ExecContext(pctx, orderID, it.ProductID, it.Title, it.SKU, it.Quantity, it.UnitPrice, line); err != nil {
			return "", err
		}
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}
	return orderID, nil
}
func (r *orderRepository) GetOrder(pctx context.Context, orderID string) (*modules.OrderDetail, error) {
	row := new(modules.Order)
	if err := r.db.GetContext(pctx, &row, `
		SELECT id, user_id, currency, total_price, status,
		       shipping_address, billing_address, created_at, updated_at
		FROM orders WHERE id=$1
	`, orderID); err != nil {
		return nil, err
	}

	items := make([]modules.OrderItemRow, 0)
	if err := r.db.SelectContext(pctx, &items, `
		SELECT id, product_id, title, sku, quantity, unit_price, total_price, created_at
		FROM order_items
		WHERE order_id=$1
		ORDER BY created_at ASC
	`, orderID); err != nil {
		return nil, err
	}

	var ship map[string]any
	var bill map[string]any
	_ = json.Unmarshal(row.Shipping, &ship)
	_ = json.Unmarshal(row.Billing, &bill)

	return &modules.OrderDetail{
		ID:           row.ID,
		UserID:       row.UserID,
		Currency:     row.Currency,
		TotalPrice:   row.TotalPrice,
		Status:       row.Status,
		ShippingAddr: ship,
		BillingAddr:  bill,
		Items:        items,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}, nil
}
func (r *orderRepository) ListOrdersByUser(pctx context.Context, userID string, limit, offset int) ([]modules.OrderListItem, error) {
	rows := make([]modules.OrderListItem, 0)
	if err := r.db.SelectContext(pctx, &rows, `
		SELECT id, total_price, status, currency, created_at
		FROM orders
		WHERE user_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset); err != nil {
		log.Printf("Error: ListOrdersByUser Failed: %v", err)
		return nil, err
	}
	return rows, nil
}
func (r *orderRepository) UpdateOrderStatus(pctx context.Context, orderID string, status string) error {
	_, err := r.db.ExecContext(pctx, `
	UPDATE orders
	SET status=$2, updated_at=CURRENT_TIMESTAMP
	WHERE id=$1
	`, orderID, status)
	if err != nil {
		log.Printf("Error: UpdateOrderStatus Failed: %v", err)
		return errors.New("update order status failed")
	}
	return nil
}

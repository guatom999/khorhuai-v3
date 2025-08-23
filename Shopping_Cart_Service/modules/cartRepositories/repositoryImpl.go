package cartRepositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/guatom999/ecommerce-shopping-cart-api/modules"
	"github.com/jmoiron/sqlx"
)

type (
	cartRepository struct {
		db *sqlx.DB
	}
)

func NewCartRepository(db *sqlx.DB) CartRepositoryInterface {
	return &cartRepository{db: db}
}

func (r *cartRepository) ResolveCartID(pctx context.Context, userID string, sessionID string) (string, error) {
	tx, err := r.db.BeginTxx(pctx, nil)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var id string
	if userID != "" {
		err = tx.GetContext(pctx, &id, `SELECT id FROM carts WHERE user_id=$1 LIMIT 1`, userID)
		if err == sql.ErrNoRows {
			err = tx.GetContext(pctx, &id,
				`INSERT INTO carts (user_id, currency) VALUES ($1,'BTH') RETURNING id`, userID)
		}
	} else if sessionID != "" {
		err = tx.GetContext(pctx, &id, `SELECT id FROM carts WHERE session_id=$1 LIMIT 1`, sessionID)
		if err == sql.ErrNoRows {
			err = tx.GetContext(pctx, &id,
				`INSERT INTO carts (session_id, currency) VALUES ($1,'BTH') RETURNING id`, sessionID)
		}
	} else {
		return "", fmt.Errorf("either user_id or session_id required")
	}
	if err != nil {
		return "", err
	}

	if err := tx.Commit(); err != nil {
		return "", err
	}

	return id, nil
}
func (r *cartRepository) UpsertCartItem(ctx context.Context, cartID, productID string, qty int, unitPrice int64, currency string) error {
	_, err := r.db.ExecContext(
		ctx,
		`
	INSERT INTO cart_items (cart_id, product_id, quantity, unit_price_cents, currency)
	VALUES ($1,$2,$3,$4,$5)
	ON CONFLICT (cart_id, product_id)
	DO UPDATE SET
	  quantity = cart_items.quantity + EXCLUDED.quantity,
	  unit_price_cents = EXCLUDED.unit_price_cents,
	  updated_at = CURRENT_TIMESTAMP;
	`,
		cartID, productID, qty, unitPrice, currency,
	)

	if err != nil {
		return err
	}
	return nil
}
func (r *cartRepository) MergeGuestToUserCart(ctx context.Context, userID string, sessionID string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var userCartID string
	err = tx.GetContext(ctx, &userCartID, `SELECT id FROM carts WHERE user_id=$1`, userID)
	if err == sql.ErrNoRows {
		if err = tx.GetContext(ctx, &userCartID,
			`INSERT INTO carts (user_id, currency) VALUES ($1,'BTH') RETURNING id`, userID); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO cart_items (cart_id, product_id, quantity, unit_price_cents, currency)
		SELECT $1, gi.product_id, gi.quantity, gi.unit_price_cents, gi.currency
		FROM cart_items gi
		JOIN carts guest_cart ON guest_cart.id = gi.cart_id
		WHERE guest_cart.session_id = $2
		ON CONFLICT (cart_id, product_id)
		DO UPDATE SET
		  quantity = cart_items.quantity + EXCLUDED.quantity,
		  updated_at = CURRENT_TIMESTAMP;
		`,
		userCartID, sessionID)
	if err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM carts WHERE session_id=$1`, sessionID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
func (r *cartRepository) GetCart(ctx context.Context, userID string, sessionID string) (*modules.CartView, error) {
	// var c modules.CartRow
	cart := new(modules.CartRow)
	var err error
	if userID != "" {
		err = r.db.GetContext(ctx, cart, `SELECT * FROM carts WHERE user_id=$1 LIMIT 1`, userID)
	} else if sessionID != "" {
		err = r.db.GetContext(ctx, cart, `SELECT * FROM carts WHERE session_id=$1 LIMIT 1`, sessionID)
	} else {
		return nil, fmt.Errorf("either user_id or session_id required")
	}
	if err == sql.ErrNoRows {
		return &modules.CartView{
				CartID:        "",
				Currency:      "BTH",
				Items:         []modules.CartItemDTO{},
				SubtotalCents: 0,
			},
			nil
	}
	if err != nil {
		return nil, err
	}

	items := make([]modules.ItemRow, 0)
	if err := r.db.SelectContext(ctx, &items, `
		SELECT id, product_id, quantity, unit_price_cents, currency, created_at, updated_at
		FROM cart_items WHERE cart_id = $1 ORDER BY created_at ASC`, cart.ID); err != nil {
		return nil, err
	}

	out := &modules.CartView{
		CartID:   cart.ID,
		Currency: cart.Currency,
		Items:    make([]modules.CartItemDTO, 0, len(items)),
	}
	var subtotal int64
	for _, it := range items {
		line := int64(it.Quantity) * it.UnitPriceCents
		subtotal += line
		out.Items = append(out.Items, modules.CartItemDTO{
			ItemID: it.ID, ProductID: it.ProductID, Quantity: it.Quantity,
			UnitPriceCents: it.UnitPriceCents, LineTotalCents: line, Currency: it.Currency,
		})
	}
	out.SubtotalCents = subtotal
	return out, nil
}

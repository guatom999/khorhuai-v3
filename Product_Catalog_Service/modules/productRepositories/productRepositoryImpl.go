package productrepositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guatom999/ecommerce-product-api/modules"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type (
	productRepository struct {
		db *sqlx.DB
	}

	productCategory struct {
		db *sqlx.DB
	}

	tempProduct struct {
		ID          uuid.UUID `db:"id"`
		Name        string    `db:"name"`
		Description *string   `db:"description"`
		Price       int64     `db:"price"`
		StockQty    int       `db:"stock_qty"`
		CreatedAt   time.Time `db:"created_at"`
		UpdatedAt   time.Time `db:"updated_at"`
		Categories  []byte    `db:"categories"`
	}
)

func NewProductRepository(db *sqlx.DB) ProductRepositoryInterface {
	return &productRepository{
		db: db,
	}
}

func NewProductCategoryRepository(db *sqlx.DB) ProductCategoryInterface {
	return &productCategory{
		db: db,
	}
}

func (r *productRepository) GetAllProductWithCategory(pctx context.Context) ([]*modules.ProductWithCategory, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	tempProducts := make([]*tempProduct, 0)

	queryString := `
	SELECT
  	  p.id, p.name, p.description, p.price, p.stock_qty,
  	  p.created_at, p.updated_at,
  	  COALESCE(
  	    json_agg(json_build_object('id', c.id, 'name', c.name, 'slug', c.slug))
  	    FILTER (WHERE c.id IS NOT NULL), '[]'
  	  ) AS categories
  	FROM products p
  	LEFT JOIN product_categories pc ON pc.product_id = p.id
  	LEFT JOIN categories c         ON c.id = pc.category_id
  	GROUP BY p.id
  	ORDER BY p.created_at DESC
  	LIMIT $1 OFFSET $2;
	`

	if err := r.db.SelectContext(ctx, &tempProducts, queryString, 10, 0); err != nil {
		log.Printf("Error: Failed to select product  %v", err)
		return make([]*modules.ProductWithCategory, 0), err
	}

	// products := make([]*modules.ProductWithCategory, 0)

	products := make([]*modules.ProductWithCategory, 0, len(tempProducts))

	for _, v := range tempProducts {
		product := &modules.ProductWithCategory{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Price:       v.Price,
			StockQty:    v.StockQty,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}

		if len(v.Categories) > 0 {
			if err := json.Unmarshal(v.Categories, &product.Categories); err != nil {
				log.Printf("Error unmarshaling categories: %v", err)

				product.Categories = []modules.Category{}
			}
		}

		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) GetAllProduct(pctx context.Context) ([]*modules.Product, error) {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	products := make([]*modules.Product, 0)

	queryString := `SELECT * FROM products`

	if err := r.db.SelectContext(ctx, &products, queryString); err != nil {
		log.Printf("Error: Failed to select Product %v", err)
		return make([]*modules.Product, 0), err
	}

	return products, nil
}

func (r *productRepository) GetOneProduct(pctx context.Context, id uuid.UUID) (*modules.Product, error) {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	product := new(modules.Product)

	queryString := `SELECT * FROM products WHERE id = $1`

	if err := r.db.GetContext(ctx, product, queryString, id); err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) CreateProduct(pctx context.Context, in *modules.CreateProductCategoryRequest) error {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var productID string
	if err = tx.GetContext(ctx, &productID, `
        INSERT INTO products (name, description, price, stock_qty)
        VALUES ($1,$2,$3,$4)
        RETURNING id;
    `, in.Name, in.Description, in.Price, in.StockQty); err != nil {
		return err
	}

	if len(in.CategorySlugs) > 0 {
		if _, err = tx.ExecContext(ctx, `
            INSERT INTO categories (name, slug)
            SELECT s.slug, s.slug
            FROM unnest($1::text[]) AS s(slug)
            ON CONFLICT (slug) DO NOTHING;
        `, pq.Array(in.CategorySlugs)); err != nil {
			return err
		}

		// 3) map product ↔ categories
		if _, err = tx.ExecContext(ctx, `
            INSERT INTO product_categories (product_id, category_id)
            SELECT $1, c.id
            FROM categories c
            WHERE c.slug = ANY($2::text[])
            ON CONFLICT (product_id, category_id) DO NOTHING;
        `, productID, pq.Array(in.CategorySlugs)); err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil

}

func (r *productRepository) UpdateProduct(pctx context.Context, id uuid.UUID, productPatch *modules.ProductPatch) error {

	fmt.Println("id is", id)

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	set := make([]string, 0)
	args := make([]any, 0)
	i := 1

	if productPatch.Name != "" {
		set = append(set, fmt.Sprintf("name = $%d", i))
		args = append(args, productPatch.Name)
		i++
	}

	if productPatch.Description != "" {
		set = append(set, fmt.Sprintf("description = $%d", i))
		args = append(args, productPatch.Description)
		i++
	}

	if productPatch.Price != 0 {
		set = append(set, fmt.Sprintf("price = $%d", i))
		args = append(args, productPatch.Price)
		i++
	}

	if productPatch.StockQty != 0 {
		set = append(set, fmt.Sprintf("stock_qty = $%d", i))
		args = append(args, productPatch.StockQty)
		i++
	}

	if len(set) == 0 {
		return errors.New("no field to update")
	}

	query := fmt.Sprintf(`
	UPDATE products
	SET %s
	WHERE id = $%d
	RETURNING id, name, description, price, stock_qty, updated_at
	`, strings.Join(set, ", "), i)

	args = append(args, id)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		log.Printf("Error: Failed to update product %v", err)
		return err
	}

	return nil
}

func (r *productRepository) DeleteProduct(pctx context.Context, id uuid.UUID) error {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	querystring := `DELETE FROM products WHERE id = $1`

	if _, err := r.db.ExecContext(ctx, querystring, id); err != nil {
		return err
	}

	deleteProductCategory := `DELETE FROM product_categories WHERE product_id = $1`

	if _, err := r.db.ExecContext(ctx, deleteProductCategory, id); err != nil {
		return err
	}

	return nil
}

func (r *productCategory) GetAllCategories(pctx context.Context) ([]*modules.Category, error) {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	categories := make([]*modules.Category, 0)

	queryString := `SELECT * FROM categories`

	if err := r.db.SelectContext(ctx, &categories, queryString); err != nil {
		return nil, err
	}

	return categories, nil
}

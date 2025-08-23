package modules

type CreateProductRequest struct {
	SKU         string   `json:"sku"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PriceCents  int64    `json:"price_cents"`
	Currency    string   `json:"currency"`
	StockQty    int      `json:"stock_qty"`
	CategoryIDs []string `json:"category_ids,omitempty"`
}

type ProductResponse struct {
	ID          string   `json:"id"`
	SKU         string   `json:"sku,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	PriceCents  int64    `json:"price_cents"`
	Currency    string   `json:"currency"`
	StockQty    int      `json:"stock_qty"`
	Categories  []string `json:"categories,omitempty"`
}

type ProductListQuery struct {
	Q        string `json:"q,omitempty"` // keyword
	Category string `json:"category,omitempty"`
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	Sort     string `json:"sort,omitempty"` // name_asc, price_desc, ...
}

type CreateProductCategoryRequest struct {
	Name          string   `json:"name" validate:"required,max=64"`
	Description   string   `json:"description,omitempty" validate:"required"`
	Price         int64    `json:"price" validate:"required,min=0"`
	StockQty      int      `json:"stock_qty" validate:"min=0"`
	CategorySlugs []string `json:"category_slugs" validate:"omitempty,dive,required,alphanum"`
}

type ProductPatch struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	StockQty    int    `json:"stock_qty"`
}

type ProductPatchReq struct {
	Id string `json:"id" validate:"required"`
	ProductPatch
}

type DeleteProductReq struct {
	Id string `json:"id" validate:"required"`
}

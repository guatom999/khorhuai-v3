package productusecases

import (
	"context"

	"github.com/guatom999/ecommerce-product-api/modules"
)

type (
	ProductUsecaseInterface interface {
		GetAllProductWithCategory(pctx context.Context) ([]*modules.ProductWithCategory, error)
		GetAllProduct(pctx context.Context) ([]*modules.Product, error)
		CreateProduct(pctx context.Context, req *modules.CreateProductCategoryRequest) error
		UpdateProduct(pctx context.Context, id string, req *modules.ProductPatchReq) error
		DeleteProduct(pctx context.Context, id string) error
	}
)

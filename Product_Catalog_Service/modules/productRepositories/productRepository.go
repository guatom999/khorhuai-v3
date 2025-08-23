package productrepositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/guatom999/ecommerce-product-api/modules"
)

type (
	ProductRepositoryInterface interface {
		GetAllProductWithCategory(pctx context.Context) ([]*modules.ProductWithCategory, error)
		CreateProduct(pctx context.Context, in *modules.CreateProductCategoryRequest) error
		UpdateProduct(pctx context.Context, id uuid.UUID, productPatch *modules.ProductPatch) error
		DeleteProduct(pctx context.Context, id uuid.UUID) error
		GetAllProduct(pctx context.Context) ([]*modules.Product, error)
		GetOneProduct(pctx context.Context, id uuid.UUID) (*modules.Product, error)
	}

	ProductCategoryInterface interface {
		GetAllCategories(pctx context.Context) ([]*modules.Category, error)
	}
)

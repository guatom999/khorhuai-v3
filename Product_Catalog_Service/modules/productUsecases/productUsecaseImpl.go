package productusecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guatom999/ecommerce-product-api/modules"
	productrepositories "github.com/guatom999/ecommerce-product-api/modules/productRepositories"
)

type (
	productUsecase struct {
		productRepo productrepositories.ProductRepositoryInterface
	}
)

func NewProductUsecase(productRepo productrepositories.ProductRepositoryInterface) ProductUsecaseInterface {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) GetAllProductWithCategory(pctx context.Context) ([]*modules.ProductWithCategory, error) {

	result, err := u.productRepo.GetAllProductWithCategory(pctx)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *productUsecase) GetAllProduct(pctx context.Context) ([]*modules.Product, error) {

	reuslt, err := u.productRepo.GetAllProduct(pctx)
	if err != nil {
		return nil, err
	}

	return reuslt, nil

}

func (u *productUsecase) UpdateProduct(pctx context.Context, id string, req *modules.ProductPatchReq) error {

	productUUid, err := uuid.Parse(req.Id)
	if err != nil {
		return errors.New("error: failed to parse id")
	}

	if err := u.productRepo.UpdateProduct(pctx, productUUid, &modules.ProductPatch{
		Name:        req.ProductPatch.Name,
		Description: req.ProductPatch.Description,
		Price:       req.ProductPatch.Price,
		StockQty:    req.ProductPatch.StockQty,
	}); err != nil {
		return err
	}

	return nil
}

func (u *productUsecase) CreateProduct(pctx context.Context, req *modules.CreateProductCategoryRequest) error {

	if err := u.productRepo.CreateProduct(pctx, req); err != nil {
		return err
	}

	return nil
}

func (u *productUsecase) DeleteProduct(pctx context.Context, id string) error {

	productUUid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("error: failed to parse id")
	}

	if err := u.productRepo.DeleteProduct(pctx, productUUid); err != nil {
		return err
	}

	return nil
}

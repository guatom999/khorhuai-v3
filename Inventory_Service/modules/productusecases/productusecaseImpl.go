package productusecases

import (
	"context"

	"github.com/guatom999/ecommerce-product-api/modules"
	"github.com/guatom999/ecommerce-product-api/modules/productrepositories"
)

type (
	productusecase struct {
		productRepo productrepositories.ProductRepositoryInterface
	}
)

func NewProductusecase(productRepo productrepositories.ProductRepositoryInterface) ProductusecaseInterface {
	return &productusecase{
		productRepo: productRepo,
	}
}

func (u *productusecase) Reserve(ctx context.Context, input modules.ReserveInput) (string, error) {

	newInput := modules.ReserveInput{
		OrderID: input.OrderID,
		UserID:  input.UserID,
		TTL:     input.TTL,
	}

	for _, item := range input.Items {
		if item.Quantity <= 0 {
			continue
		}
		newInput.Items = append(newInput.Items, item)
	}

	return u.productRepo.Reserve(ctx, newInput)
}
func (u *productusecase) Release(ctx context.Context, reservationID string) error {

	return u.productRepo.Release(ctx, reservationID)
}
func (u *productusecase) Commit(ctx context.Context, reservationID string) error {
	return u.productRepo.Commit(ctx, reservationID)
}

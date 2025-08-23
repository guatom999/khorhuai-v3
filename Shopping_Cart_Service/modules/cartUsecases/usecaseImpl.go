package cartUsecases

import (
	"context"

	"github.com/guatom999/ecommerce-shopping-cart-api/modules"
	"github.com/guatom999/ecommerce-shopping-cart-api/modules/cartRepositories"
)

type (
	cartUsecase struct {
		cartRepo cartRepositories.CartRepositoryInterface
	}
)

func NewCartUseCase(cartRepo cartRepositories.CartRepositoryInterface) CartUsecaseInterface {
	return &cartUsecase{cartRepo: cartRepo}
}
func (u *cartUsecase) GetCart(pctx context.Context, userID string, sessionID string) (*modules.CartView, error) {

	result, err := u.cartRepo.GetCart(pctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (u *cartUsecase) UpsertItem(pctx context.Context, userID string, sessionID string, productID string, qty int, unitPrice int64, currency string) (*modules.CartView, error) {
	cartID, err := u.cartRepo.ResolveCartID(pctx, userID, sessionID)
	if err != nil {
		return nil, err
	}

	if err := u.cartRepo.UpsertCartItem(pctx, cartID, productID, qty, unitPrice, currency); err != nil {
		return nil, err
	}
	return u.cartRepo.GetCart(pctx, userID, sessionID)
}
func (u *cartUsecase) Merge(pctx context.Context, userID string, sessionID string) error {

	if err := u.cartRepo.MergeGuestToUserCart(pctx, userID, sessionID); err != nil {
		return err
	}

	return nil
}

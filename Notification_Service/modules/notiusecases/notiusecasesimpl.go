package notiusecases

import (
	"context"
	"errors"
	"log"

	"github.com/guatom999/ecommerce-notification-api/modules"
	"github.com/guatom999/ecommerce-notification-api/modules/notirepositories"
)

type (
	notiusecase struct {
		notiRepo notirepositories.NotirepositoryInterface
	}
)

func NewNotiUsecase(notiRepo notirepositories.NotirepositoryInterface) NotiusecaseInterface {
	return &notiusecase{notiRepo: notiRepo}
}

func (u *notiusecase) Create(ctx context.Context, in modules.CreateInput) (string, error) {

	if in.Channel == "" || in.Recipient == "" {
		log.Printf("Errors: channel and recipient requried")
		return "", errors.New("channel and recipient requried")
	}

	notiId, err := u.notiRepo.Create(ctx, in)
	if err != nil {
		return "", err
	}

	return notiId, nil
}
func (u *notiusecase) Get(ctx context.Context, id string) (*modules.NotificationRow, error) {

	return u.notiRepo.Get(ctx, id)
}
func (u *notiusecase) List(ctx context.Context, userID, status string, limit, offset int) ([]modules.NotificationListRow, error) {

	// u.notiRepo.List(ctx, userID, status, limit, offset)

	return u.notiRepo.List(ctx, userID, status, limit, offset)
}
func (u *notiusecase) UpdateStatus(ctx context.Context, id, status string) error {
	return u.notiRepo.UpdateStatus(ctx, id, status)
}
func (u *notiusecase) AttemptSend(ctx context.Context, id string, status string, errMsg string, providerRaw map[string]any) error {

	if status != "sent" && status != "failed" {
		return errors.New("status must be sent|failed")
	}

	return u.notiRepo.AttemptSendTx(ctx, id, status, errMsg, providerRaw)
}

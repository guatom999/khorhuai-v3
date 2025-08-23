package repositories

import (
	"context"
	"time"

	"github.com/guatom999/ecommerce-user-api/config"
	"github.com/guatom999/ecommerce-user-api/modules"
)

type (
	UserRepositoryInterface interface {
		IsUserAlreadyExists(email string) bool
		CreateUser(pctx context.Context, req *modules.User) error
		EditUser(pctx context.Context, editReq *modules.EditUserReq) error
		GetUseWithCredential(pctx context.Context, email string) (*modules.User, error)
		NewAccessToken(cfg *config.Config, userId string, actExpAt time.Time) string
		NewRefreshToken(cfg *config.Config, userId string, rftExpAt time.Time) string
		CreateUserToken(pctx context.Context, token modules.UserRefreshToken) error
	}
)

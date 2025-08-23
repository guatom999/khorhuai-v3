package usecases

import (
	"context"

	"github.com/guatom999/ecommerce-user-api/modules"
)

type (
	UserUsecaseInterface interface {
		Register(pctx context.Context, req *modules.CreateUserReq) error
		Login(pctx context.Context, req *modules.LoginReq) (*modules.CustomerProfile, error)
		EditUser(pctx context.Context, req *modules.EditUserReq) error
	}
)

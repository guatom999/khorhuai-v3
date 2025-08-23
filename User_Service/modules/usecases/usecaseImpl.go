package usecases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/guatom999/ecommerce-user-api/config"
	"github.com/guatom999/ecommerce-user-api/modules"
	"github.com/guatom999/ecommerce-user-api/modules/repositories"
	"github.com/guatom999/ecommerce-user-api/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	userUsecase struct {
		userRepo repositories.UserRepositoryInterface
		cfg      *config.Config
	}
)

func NewUserUsecase(userRepo repositories.UserRepositoryInterface, cfg *config.Config) UserUsecaseInterface {
	return &userUsecase{userRepo: userRepo, cfg: cfg}
}

func (u *userUsecase) Register(pctx context.Context, req *modules.CreateUserReq) error {

	if u.userRepo.IsUserAlreadyExists(req.Email) {
		log.Printf("Error: user already exist")
		return errors.New("user already exist")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Printf("Error: failed to hashed password %v", err)
		return err
	}

	if err := u.userRepo.CreateUser(pctx, &modules.User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Name:         req.Name,
	}); err != nil {
		log.Printf("Error: failed to create user %v", err)
		return err
	}

	return nil

}

func (u *userUsecase) Login(pctx context.Context, req *modules.LoginReq) (*modules.CustomerProfile, error) {

	user, err := u.userRepo.GetUseWithCredential(pctx, req.Email)
	if err != nil {
		log.Printf("Error: failed to get user with credential %v", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("Error: failed to compare password %v", err)
		return nil, errors.New("invalid credential")
	}

	accessTokenExpireAt := utils.GetLocalTime().Add(time.Second * time.Duration(u.cfg.JWT.AccessTokenDuration))
	refreshTokenExpireAt := utils.GetLocalTime().Add(time.Second * time.Duration(u.cfg.JWT.RefreshTokenDuration))

	accessToken := u.userRepo.NewAccessToken(u.cfg, user.ID.String(), accessTokenExpireAt)
	refreshToken := u.userRepo.NewRefreshToken(u.cfg, user.ID.String(), refreshTokenExpireAt)

	if err := u.userRepo.CreateUserToken(pctx, modules.UserRefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: refreshTokenExpireAt,
	}); err != nil {
		log.Printf("Error: failed to create user token %v", err)
		return nil, err
	}

	return &modules.CustomerProfile{
		Id:           user.ID.String(),
		Email:        user.Email,
		Name:         user.Name,
		AccessToekn:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}

func (u *userUsecase) EditUser(pctx context.Context, req *modules.EditUserReq) error {

	// if err := u.userRepo.EditUser(pctx,req);

	return u.userRepo.EditUser(pctx, req)
}

package repositories

import (
	"context"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guatom999/ecommerce-user-api/config"
	"github.com/guatom999/ecommerce-user-api/modules"
	"github.com/guatom999/ecommerce-user-api/utils"
	"github.com/jmoiron/sqlx"
)

type (
	userRepository struct {
		db *sqlx.DB
	}
)

func NewUserRepository(db *sqlx.DB) UserRepositoryInterface {
	return &userRepository{db: db}
}

func (r *userRepository) IsUserAlreadyExists(email string) bool {

	queryString := `SELECT COUNT(*) FROM users WHERE email = $1`

	var count int
	err := r.db.Get(&count, queryString, email)
	if err != nil {
		return true
	}

	return count > 0
}

func (r *userRepository) CreateUser(pctx context.Context, req *modules.User) error {

	_, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	querystring := `INSERT INTO users (email, password_hash, name) VALUEs ($1 , $2 , $3)`

	_, err := r.db.Exec(querystring, req.Email, req.PasswordHash, req.Name)
	if err != nil {
		log.Printf("Error: failed to create user %v", err)
		return err
	}

	return nil
}

func (r *userRepository) GetUseWithCredential(pctx context.Context, email string) (*modules.User, error) {

	_, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	queryString := `SELECT * FROM users WHERE email = $1`

	user := new(modules.User)
	if err := r.db.Get(user, queryString, email); err != nil {
		log.Printf("Error: failed to get user with credential%v", err)
		return nil, err
	}

	return user, nil

}

func (r *userRepository) EditUser(pctx context.Context, editReq *modules.EditUserReq) error {

	_, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	queryString := `
		UPDATE users 
		SET 
			name = COALESCE($1, name),
			email = COALESCE($2, email),
			updated_at = NOW()
		WHERE id = $3
	`

	_, err := r.db.Exec(queryString, editReq.Name, editReq.Email, editReq.Id)
	if err != nil {
		log.Printf("Error: failed to edit user %v", err)
		return err
	}

	return nil
}

func (r *userRepository) CreateUserToken(pctx context.Context, token modules.UserRefreshToken) error {

	_, cancel := context.WithTimeout(pctx, time.Second*10)
	defer cancel()

	queryString := `INSERT INTO users_refresh_tokens (user_id, refresh_token, expires_at) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(queryString, token.UserID, token.Token, token.ExpiresAt)
	if err != nil {
		log.Printf("Error: failed to create user refresh token %v", err)
		return err
	}

	return nil
}

func (r *userRepository) NewAccessToken(cfg *config.Config, userId string, actExpAt time.Time) string {
	claims := modules.AuthClaims{
		Claims: &modules.Claims{
			UserId: userId,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "ecommerce-user-api.com",
			Subject:  "access_token",
			Audience: jwt.ClaimStrings{"ecommerce-user-api.com"},
			// ExpiresAt: jwt.NewNumericDate(utils.GetLocalTime().Add(time.Duration(cfg.JWT.AccessTokenDuration)))
			ExpiresAt: jwt.NewNumericDate(actExpAt),
			NotBefore: jwt.NewNumericDate(utils.GetLocalTime()),
			IssuedAt:  jwt.NewNumericDate(utils.GetLocalTime()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		log.Printf("Error: Failed to create access token %v", err)
		return ""
	}

	return accessToken

}

func (r *userRepository) NewRefreshToken(cfg *config.Config, userId string, rftExpAt time.Time) string {

	claims := modules.AuthClaims{
		Claims: &modules.Claims{
			UserId: userId,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "ecommerce-user-api.com",
			Subject:  "refresh_token",
			Audience: jwt.ClaimStrings{"ecommerce-user-api.com"},
			// ExpiresAt: jwt.NewNumericDate(utils.GetLocalTime().Add(time.Duration(cfg.JWT.RefreshTokenDuration))),
			ExpiresAt: jwt.NewNumericDate(rftExpAt),
			NotBefore: jwt.NewNumericDate(utils.GetLocalTime()),
			IssuedAt:  jwt.NewNumericDate(utils.GetLocalTime()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := token.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		log.Printf("Error: Failed to create refresh token %v", err)
		return ""
	}

	return refreshToken
}

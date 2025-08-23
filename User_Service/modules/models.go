package modules

import "github.com/golang-jwt/jwt/v5"

type (
	CreateUserReq struct {
		Email    string `json:"email" `
		Password string `json:"password" `
		Name     string `json:"name" `
	}

	LoginReq struct {
		Email    string `json:"email" `
		Password string `json:"password" `
	}

	Claims struct {
		UserId string `json:"user_id"`
	}

	AuthClaims struct {
		*Claims
		jwt.RegisteredClaims
	}

	EditUserReq struct {
		Id    string `json:"id" `
		Name  string `json:"name" `
		Email string `json:"email" `
	}

	CustomerProfile struct {
		Id           string `json:"id"`
		Email        string `json:"email"`
		Name         string `json:"name"`
		AccessToekn  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

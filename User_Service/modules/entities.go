package modules

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID           uuid.UUID `json:"id" db:"id"`
		Email        string    `json:"email" db:"email"`
		PasswordHash string    `json:"password_hash" db:"password_hash"`
		Name         string    `json:"name" db:"name"`
		CreatedAt    time.Time `json:"created_at" db:"created_at"`
		UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	}

	UserRefreshToken struct {
		ID        uuid.UUID `json:"id" db:"id"`
		UserID    uuid.UUID `json:"user_id" db:"user_id"`
		Token     string    `json:"token" db:"token"`
		Revoked   bool      `json:"revoked" db:"revoked"`
		ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
		CreatedAt time.Time `json:"created_at" db:"created_at"`
	}
)

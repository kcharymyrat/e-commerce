package data

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           uuid.UUID `db:"id" json:"id"`
	UserPhone    string    `db:"user_phone" json:"user_phone"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	IsRevoked    bool      `db:"is_revoked" json:"is_revoked"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
}

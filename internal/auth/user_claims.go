package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	Phone       string    `json:"phone"`
	FirstName   *string   `json:"first_name"`
	LastName    *string   `json:"last_name"`
	Patronomic  *string   `json:"patronomic"`
	IsActive    bool      `json:"is_active"`
	IsBanned    bool      `json:"is_banned"`
	IsStaff     bool      `json:"is_staff"`
	IsAdmin     bool      `json:"is_admin"`
	IsSuperuser bool      `json:"is_superuser"`
	jwt.RegisteredClaims
}

func newUserClaims(
	user_id uuid.UUID,
	phone string,
	firstName *string,
	lastName *string,
	patronomic *string,
	isActive bool,
	isBanned bool,
	isStaff bool,
	isAdmin bool,
	isSuperuser bool,
	duration time.Duration,
) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &UserClaims{
		UserID:      user_id,
		Phone:       phone,
		FirstName:   firstName,
		LastName:    lastName,
		Patronomic:  patronomic,
		IsActive:    isActive,
		IsBanned:    isBanned,
		IsStaff:     isStaff,
		IsAdmin:     isAdmin,
		IsSuperuser: isSuperuser,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   phone,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, err
}

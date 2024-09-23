package data

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	Patronomic   string    `json:"patronomic,omitempty"`
	DOB          time.Time `json:"dob,omitempty" validate:"birthdate"`
	Email        string    `json:"email,omitempty" validate:"email"`
}

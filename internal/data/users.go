package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID                     uuid.UUID       `json:"id" db:"id"`
	Phone                  string          `json:"phone" db:"phone"`
	PasswordHash           string          `json:"-" db:"password_hash"`                 // No need to serialize password hash
	FirstName              *string         `json:"first_name,omitempty" db:"first_name"` // Can be NULL
	LastName               *string         `json:"last_name,omitempty" db:"last_name"`   // Can be NULL
	Patronomic             *string         `json:"patronomic,omitempty" db:"patronomic"` // Can be NULL
	DOB                    *time.Time      `json:"dob,omitempty" db:"dob"`               // Can be NULL
	Email                  *string         `json:"email,omitempty" db:"email"`           // Can be NULL
	IsActive               bool            `json:"is_active" db:"is_active"`
	IsBanned               bool            `json:"is_banned" db:"is_banned"`
	IsTrusted              bool            `json:"is_trusted" db:"is_trusted"`
	InvitedById            *uuid.UUID      `json:"invited_by_id,omitempty" db:"invited_by_id"`     // Can be NULL
	InvRefId               *uuid.UUID      `json:"inv_ref_id,omitempty" db:"inv_ref_id"`           // Can be NULL
	InvProdRefId           *uuid.UUID      `json:"inv_prod_ref_id,omitempty" db:"inv_prod_ref_id"` // Can be NULL
	RefSignups             int             `json:"ref_signups" db:"ref_signups"`
	ProdRefSignups         int             `json:"prod_ref_signups" db:"prod_ref_signups"`
	ProdRefBought          int             `json:"prod_ref_bought" db:"prod_ref_bought"`
	TotalRefferals         int             `json:"total_referrals" db:"total_referrals"`
	DynamicDiscountPercent decimal.Decimal `json:"dynamic_discount_percent" db:"_dynamic_discount_percent"`
	DynDiscPercent         decimal.Decimal `json:"dyn_disc_percent" db:"dyn_disc_percent"`
	BonusPoints            decimal.Decimal `json:"bonus_points" db:"bonus_points"`
	IsStaff                bool            `json:"is_staff" db:"is_staff"`
	IsAdmin                bool            `json:"is_admin" db:"is_admin"`
	IsSuperuser            bool            `json:"is_superuser" db:"is_superuser"`
	CreatedAt              time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at" db:"updated_at"`
	CreatedById            *uuid.UUID      `json:"created_by_id,omitempty" db:"created_by_id"` // Can be NULL
	UpdatedById            *uuid.UUID      `json:"updated_by_id,omitempty" db:"updated_by_id"` // Can be NULL
}

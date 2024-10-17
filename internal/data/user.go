package data

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID                  uuid.UUID       `json:"id" db:"id" validate:"required,uuid"`
	Phone               string          `json:"phone" db:"phone" validate:"required,e164"`
	PasswordHash        string          `json:"-" db:"password_hash" validate:"required,min=8"`
	FirstName           *string         `json:"first_name,omitempty" db:"first_name" validate:"omitempty,max=50,alpha"`
	LastName            *string         `json:"last_name,omitempty"  db:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic          *string         `json:"patronomic,omitempty"  db:"patronomic" validate:"omitempty,max=50,alpha"`
	DOB                 *time.Time      `json:"dob,omitempty" db:"dob" validate:"omitempty,gte=1900-01-01"`
	Email               *string         `json:"email,omitempty" db:"email" validate:"omitempty,email"`
	IsActive            bool            `json:"is_active" db:"is_active" validate:"required"`
	IsBanned            bool            `json:"is_banned" db:"is_banned" validate:"required"`
	IsTrusted           bool            `json:"is_trusted" db:"is_trusted" validate:"required"`
	InvitedByID         *uuid.UUID      `json:"invited_by_id,omitempty" db:"invited_by_id" validate:"omitempty,uuid"`
	InvRefID            *uuid.UUID      `json:"inv_ref_id,omitempty" db:"inv_ref_id" validate:"omitempty,uuid"`
	InvProdRefID        *uuid.UUID      `json:"inv_prod_ref_id,omitempty" db:"inv_prod_ref_id" validate:"omitempty,uuid"`
	RefSignups          int             `json:"ref_signups" db:"ref_signups" validate:"gte=0"`
	ProdRefSignups      int             `json:"prod_ref_signups" db:"prod_ref_signups" validate:"gte=0"`
	ProdRefBought       int             `json:"prod_ref_bought" db:"prod_ref_bought" validate:"gte=0"`
	TotalRefferals      int             `json:"total_referrals" db:"total_referrals" validate:"gte=0"`
	WholeDynDiscPercent decimal.Decimal `json:"whole_dyn_disc_percent" db:"_dynamic_discount_percent" validate:"decimalpercent"`
	DynDiscPercent      decimal.Decimal `json:"dyn_disc_percent" db:"dyn_disc_percent" validate:"decimalpercent"`
	BonusPoints         decimal.Decimal `json:"bonus_points" db:"bonus_points" validate:"decimalgtezero"`
	IsStaff             bool            `json:"is_staff" db:"is_staff" validate:"required"`
	IsAdmin             bool            `json:"is_admin" db:"is_admin" validate:"required"`
	IsSuperuser         bool            `json:"is_superuser" db:"is_superuser" validate:"required"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt           time.Time       `json:"updated_at" db:"updated_at" validate:"required,gtefield=CreatedAt"`
	CreatedByID         *uuid.UUID      `json:"created_by_id,omitempty" db:"created_by_id" validate:"omitempty,uuid"`
	UpdatedByID         *uuid.UUID      `json:"updated_by_id,omitempty" db:"updated_by_id" validate:"omitempty,uuid"`
	Version             int             `json:"version" db:"version" validate:"required,number,min=1"`
}

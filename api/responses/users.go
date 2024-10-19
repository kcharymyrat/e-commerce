package responses

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type UserSelfResponse struct {
	ID                  uuid.UUID       `json:"id" validate:"required,uuid"`
	Phone               string          `json:"phone" validate:"required,e164"`
	FirstName           *string         `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName            *string         `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic          *string         `json:"patronomic" validate:"omitempty,max=50,alpha"`
	DOB                 *time.Time      `json:"dob" validate:"omitempty,gte=1900-01-01"`
	Email               *string         `json:"email" validate:"omitempty,email"`
	IsActive            bool            `json:"is_active" validate:"required"`
	IsBanned            bool            `json:"is_banned" validate:"required"`
	IsTrusted           bool            `json:"is_trusted" validate:"required"`
	InvitedByID         *uuid.UUID      `json:"invited_by_id" validate:"omitempty,uuid"`
	InvRefID            *uuid.UUID      `json:"inv_ref_id" validate:"omitempty,uuid"`
	InvProdRefID        *uuid.UUID      `json:"inv_prod_ref_id" validate:"omitempty,uuid"`
	RefSignups          int             `json:"ref_signups" validate:"gte=0"`
	ProdRefSignups      int             `json:"prod_ref_signups" validate:"gte=0"`
	ProdRefBought       int             `json:"prod_ref_bought" validate:"gte=0"`
	TotalRefferals      int             `json:"total_referrals" validate:"gte=0"`
	WholeDynDiscPercent decimal.Decimal `json:"whole_dyn_disc_percent" validate:"decimalpercent"`
	DynDiscPercent      decimal.Decimal `json:"dyn_disc_percent" validate:"decimalpercent"`
	BonusPoints         decimal.Decimal `json:"bonus_points" validate:"decimalgtezero"`
	IsStaff             bool            `json:"is_staff" validate:"required"`
	IsAdmin             bool            `json:"is_admin" validate:"required"`
	IsSuperuser         bool            `json:"is_superuser" validate:"required"`
	CreatedAt           time.Time       `json:"created_at" validate:"required"`
	UpdatedAt           time.Time       `json:"updated_at" validate:"required,gtefield=CreatedAt"`
	CreatedByID         *uuid.UUID      `json:"created_by_id" validate:"omitempty,uuid"`
	UpdatedByID         *uuid.UUID      `json:"updated_by_id" validate:"omitempty,uuid"`
	Version             int             `json:"version" validate:"number,min=1"`
}

type UserAdminResponse struct {
	UserSelfResponse
}

type UserPublicResponse struct {
	FirstName  *string    `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName   *string    `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic *string    `json:"patronomic" validate:"omitempty,max=50,alpha"`
	DOB        *time.Time `json:"dob" validate:"omitempty,gte=1900-01-01"`
	Email      *string    `json:"email" validate:"omitempty,email"`
	IsActive   bool       `json:"is_active" validate:"required"`
	IsStaff    bool       `json:"is_staff" validate:"required"`
}

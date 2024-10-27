package requests

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/filters"
	"github.com/shopspring/decimal"
)

// TODO: UpdateUserInput - For Superusers and Admin - IsAdmin, IsStaff, IsSuperuser fields should be required

type UsersAdminFilters struct {
	ID                      *uuid.UUID       `json:"id,omitempty" validate:"omitempty,uuid"`
	Phone                   *string          `json:"phone,omitempty" validate:"omitempty,e164"`
	Email                   *string          `json:"email,omitempty" validate:"omitempty,email"`
	IsActice                *bool            `json:"is_active,omitempty" validate:"omitempty"`
	IsBanned                *bool            `json:"is_banned,omitempty" validate:"omitempty"`
	IsTrusted               *bool            `json:"is_trusted,omitempty" validate:"omitempty"`
	IsInvited               *bool            `json:"is_invited,omitempty" validate:"omitempty"`
	RefSignupsFrom          *int             `json:"ref_signups_from" validate:"omitempty,min=0"`
	RefSignupsTo            *int             `json:"ref_signups_to" validate:"omitempty,min=0,gtefield=RefSignupsFrom"`
	ProdRefSignupsFrom      *int             `json:"p_ref_signups_from" validate:"omitempty,min=0"`
	ProdRefSignupsTo        *int             `json:"p_ref_signups_to" validate:"omitempty,min=0,gtefield=ProdRefSignupsFrom"`
	ProdRefBoughtFrom       *int             `json:"p_ref_bought_from" validate:"omitempty,min=0"`
	ProdRefBoughtTo         *int             `json:"p_ref_bought_to" validate:"omitempty,min=0,gtefield=ProdRefBoughtFrom"`
	WholeDynDiscPercentFrom *decimal.Decimal `json:"whole_ddp_from" validate:"omitempty,decimalpercent"`
	WholeDynDiscPercentTo   *decimal.Decimal `json:"whole_ddp_to" validate:"omitempty,decimalpercent,gtefield=WholeDynDiscPercentFrom"`
	DynDiscPercentFrom      *decimal.Decimal `json:"ddp_from" validate:"omitempty,decimalpercent"`
	DynDiscPercentTo        *decimal.Decimal `json:"ddp_to" validate:"omitempty,decimalpercent,gtefield=DynDiscPercentFrom"`
	BonusPointsFrom         *decimal.Decimal `json:"bonus_from" validate:"omitempty,decimalgtezero"`
	BonusPointsTo           *decimal.Decimal `json:"bonus_to" validate:"omitempty,decimalgtezero,gtefield=BonusPointsFrom"`
	IsStaff                 *bool            `json:"is_staff,omitempty" validate:"omitempty"`
	IsAdmin                 *bool            `json:"is_admin,omitempty" validate:"omitempty"`
	IsSuperuser             *bool            `json:"is_superuser,omitempty" validate:"omitempty"`
	filters.SearchFilter
	filters.CreatedUpdatedAtFilter
	filters.CreatedUpdatedByFilter
	filters.SortListFilter
	filters.PaginationFilter
}

type UserAdminCreate struct {
	Phone       string     `json:"phone" validate:"required,e164"`
	Password    string     `json:"password" validate:"required,min=8,max=72,password"`
	FirstName   *string    `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName    *string    `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic  *string    `json:"patronomic" validate:"omitempty,max=50,alpha"`
	Email       *string    `json:"email" validate:"omitempty,email"`
	IsActive    bool       `json:"is_active" validate:"required"`
	CreatedByID *uuid.UUID `json:"created_by_id,omitempty" validate:"omitempty,uuid"`
}

type UserAdminUpdate struct {
	Phone       string    `json:"phone" validate:"required,e164"`
	Password    string    `json:"password" validate:"required,min=8,max=72,password"`
	FirstName   *string   `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName    *string   `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic  *string   `json:"patronomic" validate:"omitempty,max=50,alpha"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	IsActive    bool      `json:"is_active" validate:"required"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"uuid"`
}

type UserAdminPartialUpdate struct {
	Phone       *string   `json:"phone" validate:"required,e164"`
	Password    *string   `json:"password" validate:"required,min=8,max=72,password"`
	FirstName   *string   `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName    *string   `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic  *string   `json:"patronomic" validate:"omitempty,max=50,alpha"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	IsActive    *bool     `json:"is_active" validate:"required"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"uuid"`
}

type UserSelfUpdate struct {
	Phone       string    `json:"phone" validate:"required,e164"`
	Password    string    `json:"password" validate:"required,min=8,max=72,password"`
	FirstName   *string   `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName    *string   `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic  *string   `json:"patronomic" validate:"omitempty,max=50,alpha"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"uuid"`
}

type UserSelfPartialUpdate struct {
	Phone       *string   `json:"phone" validate:"required,e164"`
	Password    *string   `json:"password" validate:"required,min=8,max=72,password"`
	FirstName   *string   `json:"first_name" validate:"omitempty,max=50,alpha"`
	LastName    *string   `json:"last_name" validate:"omitempty,max=50,alpha"`
	Patronomic  *string   `json:"patronomic" validate:"omitempty,max=50,alpha"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"uuid"`
}

type AdminUserLoginReq struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=8,max=72,password"`
}

type UserLoginReq struct {
	Phone string `json:"phone" validate:"required,e164"`
}

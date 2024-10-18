package requests

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/filters"
)

type LanguagesAdminFilters struct {
	filters.PaginationFilter
}

type LanguageAdminCreate struct {
	Code        string    `json:"code" validate:"required,min=2,max=10"`
	Name        string    `json:"name" validate:"required,min=2,max=50"`
	CreatedByID uuid.UUID `json:"created_by_id" validate:"required,uuid"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

type LanguageAdminUpdate struct {
	Code        string    `json:"code" validate:"required,min=2,max=10"`
	Name        string    `json:"name" validate:"required,min=2,max=50"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

type LanguageAdminPartialUpdate struct {
	Code        *string   `json:"code,omitempty" validate:"omitempty,min=2,max=10"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

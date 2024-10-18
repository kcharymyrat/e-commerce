package requests

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/filters"
)

type ListLanguagesFilters struct {
	filters.PaginationFilter
}

type CreateLanguageInput struct {
	Code        string    `json:"code" validate:"required,min=2,max=10"`
	Name        string    `json:"name" validate:"required,min=2,max=50"`
	CreatedByID uuid.UUID `json:"created_by_id" validate:"required,uuid"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

type UpdateLanguageInput struct {
	Code        string    `json:"code" validate:"required,min=2,max=10"`
	Name        string    `json:"name" validate:"required,min=2,max=50"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

type PartialUpdateLanguageInput struct {
	Code        *string   `json:"code,omitempty" validate:"omitempty,min=2,max=10"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=2,max=50"`
	UpdatedByID uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

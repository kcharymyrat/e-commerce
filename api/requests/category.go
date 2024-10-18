package requests

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/filters"
)

type ListCategoriesFilters struct {
	Names     []string    `json:"names" validate:"omitempty,dive,max=50"`
	Slugs     []string    `json:"slugs" validate:"omitempty,dive,max=50,slug"`
	ParentIDs []uuid.UUID `json:"parent_ids" validate:"omitempty,dive,uuid"`
	filters.SearchFilter
	filters.CreatedUpdatedAtFilter
	filters.CreatedUpdatedByFilter
	filters.SortListFilter
	filters.PaginationFilter
}

type CreateCategoryInput struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty" validate:"omitempty,uuid"`
	Name        string     `json:"name" validate:"required,min=3,max=50"`
	Slug        string     `json:"slug" validate:"required,slug"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=500"`
	ImageUrl    string     `json:"image_url" validate:"required,url"`
	CreatedByID uuid.UUID  `json:"created_by_id" validate:"required,uuid"`
	UpdatedByID uuid.UUID  `json:"updated_by_id" validate:"required,uuid"`
}

type UpdateCategoryInput struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty" validate:"omitempty,uuid"`
	Name        string     `json:"name" validate:"required,min=3,max=50"`
	Slug        string     `json:"slug" validate:"required,slug"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=500"`
	ImageUrl    string     `json:"image_url" validate:"required,url"`
	UpdatedByID uuid.UUID  `json:"updated_by_id" validate:"required,uuid"`
}

type PartialUpdateCategoryInput struct {
	ParentID    *uuid.UUID `json:"parent_id,omitempty"  validate:"omitempty,uuid"`
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=3,max=50"`
	Slug        *string    `json:"slug,omitempty" validate:"omitempty,slug"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=500"`
	ImageUrl    *string    `json:"image_url,omitempty" validate:"omitempty,url"`
	UpdatedByID uuid.UUID  `json:"updated_by_id" validate:"required,uuid"`
}

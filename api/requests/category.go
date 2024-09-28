package requests

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/filters"
)

type ListCategoriesInput struct {
	Names     []string    `json:"names" validate:"omitempty,dive,max=50"`
	Slugs     []string    `json:"slugs" validate:"omitempty,dive,max=50,slug"`
	ParentIDs []uuid.UUID `json:"parent_ids" validate:"omitempty,dive,uuid"`
	filters.SearchFilters
	filters.CreatedUpdatedAtFilter
	filters.CreatedUpdatedByFilters
	filters.SortListFilters
	filters.PaginationFilters
}

type CreateCategoryInput struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	ImageUrl    string     `json:"image_url"`
	Description *string    `json:"description"`
	CreatedByID uuid.UUID  `json:"created_by_id"`
	UpdatedByID uuid.UUID  `json:"updated_by_id"`
}

type UpdateCategoryInput struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	ImageUrl    string     `json:"image_url"`
	Description *string    `json:"description"`
	UpdatedByID uuid.UUID  `json:"updated_by_id"`
}

type PartialUpdateCategoryInput struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        *string    `json:"name"`
	Slug        *string    `json:"slug"`
	ImageUrl    *string    `json:"image_url"`
	Description *string    `json:"description"`
	UpdatedByID uuid.UUID  `json:"updated_by_id"`
}

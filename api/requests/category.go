package requests

import "github.com/google/uuid"

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
	UpdatedByID uuid.UUID  `json:"created_by_id"`
}

type PartialUpdateCategoryInput struct {
	ParentID    *uuid.UUID `json:"parent_id"`
	Name        *string    `json:"name"`
	Slug        *string    `json:"slug"`
	ImageUrl    *string    `json:"image_url"`
	Description *string    `json:"description"`
	UpdatedByID uuid.UUID  `json:"created_by_id"`
}

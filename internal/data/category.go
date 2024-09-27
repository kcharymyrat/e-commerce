package data

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID  `json:"id" db:"id" validate:"required,uuid"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id" validate:"omitempty,uuid"`
	Name        string     `json:"name" db:"name" validate:"required,min=3,max=50"`
	Slug        string     `json:"slug" db:"slug" validate:"required,slug"`
	Description *string    `json:"description,omitempty" db:"description" validate:"omitempty,max=500"`
	ImageUrl    string     `json:"image_url" db:"image_url" validate:"required,url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at" validate:"required"`
	CreatedByID uuid.UUID  `json:"created_by_id" db:"created_by_id" validate:"required,uuid"`
	UpdatedByID uuid.UUID  `json:"updated_by_id" db:"updated_by_id" validate:"required,uuid"`
	Version     int        `json:"version" db:"version" validate:"gte=0"`
}

// func ValidateCategory(v *validator.Validator, category *Category) {
// 	v.Check(category.Name != "", "name", "must be provided")
// 	v.Check(category.Slug != "", "slug", "must be provided")
// 	v.Check(category.ImageUrl != "", "image_url", "must be provided")
// 	v.Check(category.CreatedByID != uuid.Nil, "created_by_id", "must be provided")
// 	v.Check(category.UpdatedByID != uuid.Nil, "updated_by_id", "must be provided")

// 	v.Check(len([]byte(category.Name)) <= 50, "name", "must not be more than 50 bytes long")
// 	v.Check(len([]byte(category.Slug)) <= 50, "slug", "must not be more than 50 bytes long")

// 	v.Check(category.UpdatedAt.After(time.Now()), "updated_at", "can not be in the future")
// 	v.Check(category.CreatedAt.After(category.UpdatedAt), "created_at", "can not be later than updated_at")
// }

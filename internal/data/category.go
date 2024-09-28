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

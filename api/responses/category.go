package responses

import (
	"time"

	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

type CategoryPublicResponse struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" db:"parent_id"`
	Name        string     `json:"name" db:"name"`
	Slug        string     `json:"slug" db:"slug"`
	Description *string    `json:"description,omitempty" db:"description"`
	ImageUrl    string     `json:"image_url" db:"image_url"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CategoryAdminResponse struct {
	data.Category
}

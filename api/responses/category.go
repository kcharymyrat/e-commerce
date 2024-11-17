package responses

import (
	"time"

	"github.com/google/uuid"
)

type CategoryWithTranslationsAdminResponse struct {
	Category     CategoryAdminResponse        `json:"category"`
	Translations map[string]map[string]string `json:"translations"`
}

type CategoryWithTranslationsPublicResponse struct {
	Category     CategoryPublicResponse       `json:"category"`
	Translations map[string]map[string]string `json:"translations"`
}

type CategoryAdminResponse struct {
	ID          uuid.UUID  `json:"id"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description *string    `json:"description,omitempty"`
	ImageUrl    string     `json:"image_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CreatedByID uuid.UUID  `json:"created_by_id"`
	UpdatedByID uuid.UUID  `json:"updated_by_id"`
	Version     int        `json:"version"`
}

type CategoryPublicResponse struct {
	ID          uuid.UUID  `json:"id" format:"uuid"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty" format:"uuid"`
	Name        string     `json:"name" example:"Electronics"`
	Slug        string     `json:"slug" format:"slug" example:"electronics"`
	Description *string    `json:"description,omitempty"`
	ImageUrl    string     `json:"image_url" example:"https://example.com/image.jpg" format:"url"`
	CreatedAt   time.Time  `json:"created_at" format:"date-time"`
	UpdatedAt   time.Time  `json:"updated_at" format:"date-time"`
}

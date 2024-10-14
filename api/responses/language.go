package responses

import (
	"time"

	"github.com/google/uuid"
)

type LanguageManagerResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedByID uuid.UUID `json:"created_by_id"`
	UpdatedByID uuid.UUID `json:"updated_by_id"`
	Version     int       `json:"version"`
}

type LanguagePublicResponse struct {
	ID   uuid.UUID `json:"id"`
	Code string    `json:"code"`
	Name string    `json:"name"`
}

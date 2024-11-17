package data

import (
	"time"

	"github.com/google/uuid"
)

type Language struct {
	ID          uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	Code        string    `json:"code" db:"code" validate:"required,min=2,max=10"`
	Name        string    `json:"name" db:"name" validate:"required,min=2,max=50"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at" validate:"required,gtefield=CreatedAt"`
	CreatedByID uuid.UUID `json:"created_by_id" db:"created_by_id" validate:"required,uuid"`
	UpdatedByID uuid.UUID `json:"updated_by_id" db:"updated_by_id" validate:"required,uuid"`
	Version     int       `json:"version" db:"version" validate:"required,number,min=1"`
}

type LanguageWithTranslations struct {
	Language     *Language
	Translations []*Translation
}

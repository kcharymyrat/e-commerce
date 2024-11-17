package responses

import (
	"time"

	"github.com/google/uuid"
)

type TranslationAdminResponse struct {
	ID                  uuid.UUID `json:"id" validate:"uuid"`
	LanguageCode        string    `json:"language_code" validate:"min=2,max=10"`
	EntityID            uuid.UUID `json:"entity_id" validate:"uuid"`
	TableName           string    `json:"table_name" validate:"min=1,max=50"`
	FieldName           string    `json:"field_name" validate:"min=1,max=50"`
	TranslatedFieldName string    `json:"translated_field_name" validate:"min=1,max=50"`
	TranslatedValue     string    `json:"translated_value" validate:"min=1"`
	CreatedAt           time.Time `json:"created_at" validate:"required"`
	UpdatedAt           time.Time `json:"updated_at" validate:"required,gtefield=CreatedAt"`
	CreatedByID         uuid.UUID `json:"created_by_id" validate:"required,uuid"`
	UpdatedByID         uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
	Version             int       `json:"version" db:"version" validate:"required,number,min=1"`
}

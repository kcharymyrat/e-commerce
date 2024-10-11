package data

import (
	"time"

	"github.com/google/uuid"
)

type Translation struct {
	ID              uuid.UUID `json:"id" db:"id" validate:"uuid"`
	LanguageCode    string    `json:"language_code" db:"language_code" validate:"min=2,max=10"`
	EntityID        uuid.UUID `json:"entity_id" db:"entity_id" validate:"uuid"`
	TableName       string    `json:"table_name" db:"table_name" validate:"min=1"`
	FieldName       string    `json:"field_name" db:"field_name" validate:"min=1"`
	TranslatedValue string    `json:"translated_value" db:"translated_value" validate:"min=1"`
	CreatedAt       time.Time `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at" validate:"required,gtefield=CreatedAt"`
	CreatedByID     uuid.UUID `json:"created_by_id" db:"created_by_id" validate:"required,uuid"`
	UpdatedByID     uuid.UUID `json:"updated_by_id" db:"updated_by_id" validate:"required,uuid"`
	Version         int       `json:"version" db:"version" validate:"required,number,min=1"`
}

// CREATE TABLE IF NOT EXISTS translations (
//     id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
//     language_code varchar(10) NOT NULL,
//     entity_id uuid NOT NULL,
//     table_name varchar(50) NOT NULL,
//     field_name varchar(50) NOT NULL,
//     translated_value text NOT NULL,

//     created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
//     created_by_id uuid NOT NULL,
//     updated_by_id uuid NOT NULL,

//     version integer NOT NULL DEFAULT 1,

//     CHECK (updated_at >= created_at),
//     UNIQUE (entity_id, language_code, table_name, field_name)
// );

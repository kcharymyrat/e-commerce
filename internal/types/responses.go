package types

import (
	"github.com/google/uuid"
)

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type TranslationResponse struct {
	ID                  uuid.UUID `json:"id" format:"uuid"`
	LanguageCode        string    `json:"language_code" example:"ru"`
	EntityID            uuid.UUID `json:"entity_id" format:"uuid"`
	TableName           string    `json:"table_name" example:"products"`
	FieldName           string    `json:"field_name" example:"name"`
	TranslatedFieldName string    `json:"translated_field_name" example:"название"`
	TranslatedValue     string    `json:"translated_value" example:"ноутбук"`
}

type DetailResponse[T any] struct {
	Data         *T                     `json:"data"`
	Translations []*TranslationResponse `json:"translations"`
}

type PaginatedResponse[T any] struct {
	Metadata PaginationMetadata   `json:"metadata"`
	Results  []*DetailResponse[T] `json:"results"`
}

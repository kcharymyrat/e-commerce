package requests

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/filters"
)

type ListTranslationsInput struct {
	LanguageCodes []string    `json:"language_codes" validate:"omitempty,dive,max=10"`
	TableNames    []string    `json:"table_names" validate:"omitempty,dive,max=255"`
	FieldNames    []string    `json:"field_names" validate:"omitempty,dive,max=255"`
	EntityIDs     []uuid.UUID `json:"entity_ids" validate:"omitempty,dive,uuid"`
	filters.SearchFilters
	filters.CreatedUpdatedAtFilter
	filters.CreatedUpdatedByFilters
	filters.SortListFilters
	filters.PaginationFilters
}

type CreateTranslationInput struct {
	LanguageCode        string    `json:"language_code" validate:"min=2,max=10"`
	EntityID            uuid.UUID `json:"entity_id" validate:"uuid"`
	TableName           string    `json:"table_name" validate:"min=1,max=50"`
	FieldName           string    `json:"field_name" validate:"min=1,max=50"`
	TranslatedFieldName string    `json:"translated_field_name" validate:"min=1,max=50"`
	TranslatedValue     string    `json:"translated_value" validate:"min=1"`
	CreatedByID         uuid.UUID `json:"created_by_id" validate:"required,uuid"`
	UpdatedByID         uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

type UpdateTranslationInput struct {
	LanguageCode        string    `json:"language_code" validate:"min=2,max=10"`
	EntityID            uuid.UUID `json:"entity_id" validate:"uuid"`
	TableName           string    `json:"table_name" validate:"min=1,max=50"`
	FieldName           string    `json:"field_name" validate:"min=1,max=50"`
	TranslatedFieldName string    `json:"translated_field_name" validate:"min=1,max=50"`
	TranslatedValue     string    `json:"translated_value" validate:"min=1"`
	UpdatedByID         uuid.UUID `json:"updated_by_id" validate:"required,uuid"`
}

type PartialUpdateTranslationInput struct {
	LanguageCode        *string    `json:"language_code,omitempty" validate:"omitempty,min=2,max=10"`
	EntityID            *uuid.UUID `json:"entity_id,omitempty" validate:"omitempty,uuid"`
	TableName           *string    `json:"table_name,omitempty" validate:"omitempty,min=1,max=50"`
	FieldName           *string    `json:"field_name,omitempty" validate:"omitempty,min=1,max=50"`
	TranslatedFieldName *string    `json:"translated_field_name" validate:"omitempty,min=1,max=50"`
	TranslatedValue     *string    `json:"translated_value,omitempty" validate:"omitempty,min=1"`
	UpdatedByID         uuid.UUID  `json:"updated_by_id" validate:"required,uuid"`
}

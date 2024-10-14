package mappers

import (
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateTranslationInputToTranslationMapper(input *requests.CreateTranslationInput) *data.Translation {
	return &data.Translation{
		LanguageCode:    input.LanguageCode,
		EntityID:        input.EntityID,
		TableName:       input.TableName,
		FieldName:       input.FieldName,
		TranslatedValue: input.TranslatedValue,
		CreatedByID:     input.CreatedByID,
		UpdatedByID:     input.UpdatedByID,
	}
}

func TranslationToTranslationManagerResponseMappper(tr *data.Translation) *responses.TranslationManagerResponse {
	return &responses.TranslationManagerResponse{
		ID:              tr.ID,
		LanguageCode:    tr.LanguageCode,
		EntityID:        tr.EntityID,
		TableName:       tr.TableName,
		FieldName:       tr.FieldName,
		TranslatedValue: tr.TranslatedValue,
	}
}

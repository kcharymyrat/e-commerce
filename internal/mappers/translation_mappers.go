package mappers

import (
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateTranslationInputToTranslationMapper(input *requests.TranslationAdminCreate) *data.Translation {
	return &data.Translation{
		LanguageCode:        input.LanguageCode,
		EntityID:            input.EntityID,
		TableName:           input.TableName,
		FieldName:           input.FieldName,
		TranslatedFieldName: input.TranslatedFieldName,
		TranslatedValue:     input.TranslatedValue,
		CreatedByID:         input.CreatedByID,
		UpdatedByID:         input.UpdatedByID,
	}
}

func TranslationToTranslationManagerResponseMappper(tr *data.Translation) *responses.TranslationAdminResponse {
	return &responses.TranslationAdminResponse{
		ID:                  tr.ID,
		LanguageCode:        tr.LanguageCode,
		EntityID:            tr.EntityID,
		TableName:           tr.TableName,
		FieldName:           tr.FieldName,
		TranslatedFieldName: tr.TranslatedFieldName,
		TranslatedValue:     tr.TranslatedValue,
		CreatedAt:           tr.CreatedAt,
		UpdatedAt:           tr.UpdatedAt,
		CreatedByID:         tr.CreatedByID,
		UpdatedByID:         tr.UpdatedByID,
		Version:             tr.Version,
	}
}

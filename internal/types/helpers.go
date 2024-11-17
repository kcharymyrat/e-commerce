package types

import "github.com/kcharymyrat/e-commerce/internal/data"

func NewTranslationResponse(data *data.Translation) *TranslationResponse {
	return &TranslationResponse{
		ID:                  data.ID,
		LanguageCode:        data.LanguageCode,
		EntityID:            data.EntityID,
		TableName:           data.TableName,
		FieldName:           data.FieldName,
		TranslatedFieldName: data.TranslatedFieldName,
		TranslatedValue:     data.TranslatedValue,
	}

}

func NewDetailResponse[T any](data *T, trData []*data.Translation) *DetailResponse[T] {

	translations := []*TranslationResponse{}
	for _, tr := range trData {
		translations = append(translations, NewTranslationResponse(tr))
	}

	return &DetailResponse[T]{
		Data:         data,
		Translations: translations,
	}
}

func NewPaginatedResponse[T any](data []*T, trData []*data.Translation) *PaginatedResponse[T] {
	return &PaginatedResponse[T]{}
}

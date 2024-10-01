package mappers

import (
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateLanguageInputToLanguageMapper(input *requests.CreateLanguageInput) *data.Language {
	return &data.Language{
		Name:        input.Name,
		Code:        input.Code,
		CreatedByID: input.CreatedByID,
		UpdatedByID: input.UpdatedByID,
	}
}

func LanguageToLanguageManagerResponseMapper(input *data.Language) *responses.LanguageManagerResponse {
	return &responses.LanguageManagerResponse{
		ID:          input.ID,
		Code:        input.Code,
		Name:        input.Name,
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
		CreatedByID: input.CreatedByID,
		UpdatedByID: input.UpdatedByID,
	}
}

func LanguageToLanguagePublicResponseMapper(input *data.Language) *responses.LanguagePublicResponse {
	return &responses.LanguagePublicResponse{
		ID:   input.ID,
		Code: input.Code,
		Name: input.Name,
	}
}

package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/types"
)

func ListLanguagesPublicService(
	app *app.Application,
	filters *requests.LanguagesAdminFilters,
	langCode string,
) ([]*data.LanguageWithTranslations, types.PaginationMetadata, error) {

	langs, metadata, err := app.Repositories.Languages.List(filters)
	if err != nil {
		return nil, types.PaginationMetadata{}, err
	}

	langsWithTrans := make([]*data.LanguageWithTranslations, len(langs))
	for _, lang := range langs {
		fieldsToTranslate := []string{"name", "description"}
		translations, err := GetTranslationSlice(app, lang.ID, langCode, fieldsToTranslate)
		if err != nil {
			return nil, types.PaginationMetadata{}, err
		}

		result := &data.LanguageWithTranslations{
			Language:     lang,
			Translations: translations,
		}

		langsWithTrans = append(langsWithTrans, result)
	}
	return langsWithTrans, metadata, nil
}

func GetLanguagePublicService(
	app *app.Application, id uuid.UUID, langCode string,
) (*data.LanguageWithTranslations, error) {
	lang, err := app.Repositories.Languages.GetByID(id)
	if err != nil {
		return nil, err
	}

	fieldsToTranslate := []string{"name"}
	translations, err := GetTranslationSlice(app, lang.ID, langCode, fieldsToTranslate)
	if err != nil {
		return nil, err
	}

	result := &data.LanguageWithTranslations{
		Language:     lang,
		Translations: translations,
	}
	return result, nil
}

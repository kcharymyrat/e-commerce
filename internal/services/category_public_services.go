package services

import (
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/types"
)

func ListCategoriesPublicService(
	app *app.Application,
	filters *requests.CategoriesAdminFilters,
	langCode string,
) ([]*data.CategoryWithTranslations, types.PaginationMetadata, error) {

	categories, metadata, err := app.Repositories.Categories.List(filters)
	if err != nil {
		return nil, types.PaginationMetadata{}, err
	}

	catsWithTrans := make([]*data.CategoryWithTranslations, len(categories))
	for _, category := range categories {
		fieldsToTranslate := []string{"name", "description"}
		translations, err := GetTranslationSlice(app, category.ID, langCode, fieldsToTranslate)
		if err != nil {
			return nil, types.PaginationMetadata{}, err
		}

		result := &data.CategoryWithTranslations{
			Category:     category,
			Translations: translations,
		}

		catsWithTrans = append(catsWithTrans, result)
	}
	return catsWithTrans, metadata, nil
}

func GetCategoryBySlugPublicService(
	app *app.Application, slug string, langCode string,
) (*data.CategoryWithTranslations, error) {
	category, err := GetCategoryBySlugService(app, slug)
	if err != nil {
		return nil, err
	}

	fieldsToTranslate := []string{"name", "description"}
	translations, err := GetTranslationSlice(app, category.ID, langCode, fieldsToTranslate)
	if err != nil {
		return nil, err
	}

	result := &data.CategoryWithTranslations{
		Category:     category,
		Translations: translations,
	}
	return result, nil
}

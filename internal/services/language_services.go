package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateLanguageService(app *app.Application, language *data.Language) error {
	return app.Repositories.Languages.Create(language)
}

func GetLanguageService(app *app.Application, id uuid.UUID) (*data.Language, error) {
	return app.Repositories.Languages.GetByID(id)
}

func ListLanguagesService(
	app *app.Application,
	filters *requests.ListLanguagesFilters,
) ([]*data.Language, common.Metadata, error) {
	return app.Repositories.Languages.List(filters)
}

func UpdateLanguageService(
	app *app.Application,
	input *requests.UpdateLanguageInput,
	language *data.Language,
) error {
	language.Name = input.Name
	language.Code = input.Code
	language.UpdatedByID = input.UpdatedByID

	return app.Repositories.Languages.Update(language)
}

func PartialUpdateLanguageService(
	app *app.Application,
	input *requests.PartialUpdateLanguageInput,
	language *data.Language,
) error {
	if input.Name != nil {
		language.Name = *input.Name
	}
	if input.Code != nil {
		language.Code = *input.Code
	}
	language.UpdatedByID = input.UpdatedByID

	return app.Repositories.Languages.Update(language)
}

func DeleteLanguageService(app *app.Application, id uuid.UUID) error {
	return app.Repositories.Languages.Delete(id)
}

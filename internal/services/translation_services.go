package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateTranslationService(app *app.Application, tr *data.Translation) error {
	return app.Repositories.Translations.Create(tr)
}

func GetTranslationService(app *app.Application, id uuid.UUID) (*data.Translation, error) {
	return app.Repositories.Translations.GetByID(id)
}

func ListTranslationsService(
	app *app.Application,
	filters *requests.TranslationsAdminFilters,
) ([]*data.Translation, common.Metadata, error) {
	return app.Repositories.Translations.List(filters)
}

func UpdateTranslationService(
	app *app.Application,
	input *requests.TranslationAdminUpdate,
	tr *data.Translation,
) error {
	tr.LanguageCode = input.LanguageCode
	tr.EntityID = input.EntityID
	tr.TableName = input.TableName
	tr.FieldName = input.FieldName
	tr.TranslatedFieldName = input.TranslatedFieldName
	tr.TranslatedValue = input.TranslatedValue
	tr.UpdatedByID = input.UpdatedByID

	return app.Repositories.Translations.Update(tr)
}

func PartialUpdateTranslationService(
	app *app.Application,
	input *requests.TranslationAdminPartialUpdate,
	tr *data.Translation,
) error {
	if input.LanguageCode != nil {
		tr.LanguageCode = *input.LanguageCode
	}
	if input.EntityID != nil {
		tr.EntityID = *input.EntityID
	}
	if input.TableName != nil {
		tr.TableName = *input.TableName
	}
	if input.FieldName != nil {
		tr.FieldName = *input.FieldName
	}
	if input.TranslatedFieldName != nil {
		tr.TranslatedFieldName = *input.TranslatedFieldName
	}
	if input.TranslatedValue != nil {
		tr.TranslatedValue = *input.TranslatedValue
	}
	tr.UpdatedByID = input.UpdatedByID

	return app.Repositories.Translations.Update(tr)
}

func DeleteTranslationService(app *app.Application, id uuid.UUID) error {
	return app.Repositories.Translations.Delete(id)
}

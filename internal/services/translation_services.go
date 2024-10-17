package services

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateTranslationService(app *app.Application, tr *data.Translation) error {
	err := app.Repositories.Translations.Create(tr)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return common.TransformPgErrToCustomError(pgErr)
		}
		return err
	}
	return nil
}

func GetTranslationService(app *app.Application, id uuid.UUID) (*data.Translation, error) {
	return app.Repositories.Translations.GetByID(id)
}

func ListTranslationsService(
	app *app.Application,
	listTr *requests.ListTranslationsInput,
) ([]*data.Translation, common.Metadata, error) {
	return app.Repositories.Translations.List(
		listTr.LanguageCodes,
		listTr.TableNames,
		listTr.FieldNames,
		listTr.EntityIDs,
		listTr.Search,
		listTr.CreatedAtFrom,
		listTr.CreatedAtUpTo,
		listTr.UpdatedAtFrom,
		listTr.UpdatedAtUpTo,
		listTr.CreatedByIDs,
		listTr.UpdatedByIDs,
		listTr.Sorts,
		listTr.SortSafeList,
		listTr.Page,
		listTr.PageSize,
	)
}

func UpdateTranslationService(
	app *app.Application,
	input *requests.UpdateTranslationInput,
	tr *data.Translation,
) error {
	tr.LanguageCode = input.LanguageCode
	tr.EntityID = input.EntityID
	tr.TableName = input.TableName
	tr.FieldName = input.FieldName
	tr.TranslatedFieldName = input.TranslatedFieldName
	tr.TranslatedValue = input.TranslatedValue
	tr.UpdatedByID = input.UpdatedByID

	err := app.Repositories.Translations.Update(tr)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return common.TransformPgErrToCustomError(pgErr)
		}
		return err
	}
	return nil
}

func PartialUpdateTranslationService(
	app *app.Application,
	input *requests.PartialUpdateTranslationInput,
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

	err := app.Repositories.Translations.Update(tr)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return common.TransformPgErrToCustomError(pgErr)
		}
		return err
	}
	return nil
}

func DeleteTranslationService(app *app.Application, id uuid.UUID) error {
	return app.Repositories.Translations.Delete(id)
}

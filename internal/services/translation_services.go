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

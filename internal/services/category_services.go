package services

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func ListCategoriesService(
	app *app.Application,
	input requests.ListCategoriesInput,
) ([]*data.Category, common.Metadata, error) {
	return app.Repositories.Categories.GetAll(
		input.Names,
		input.Slugs,
		input.ParentIDs,
		input.Search,
		input.CreatedAtFrom,
		input.CreatedAtUpTo,
		input.UpdatedAtFrom,
		input.UpdatedAtUpTo,
		input.CreatedByIDs,
		input.UpdatedByIDs,
		input.Sorts,
		input.SortSafeList,
		input.Page,
		input.PageSize,
	)

}

func CreateCategoryService(app *app.Application, category *data.Category) error {
	err := app.Repositories.Categories.Insert(category)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return common.TransformPgErrToCustomError(pgErr)
		}
		return err
	}
	return nil
}

func GetCategoryService(app *app.Application, id uuid.UUID) (*data.Category, error) {
	return app.Repositories.Categories.Get(id)
}

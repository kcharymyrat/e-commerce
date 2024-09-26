package services

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

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

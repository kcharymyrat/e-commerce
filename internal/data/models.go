package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Categories CategoryModel
}

func NewModels(dbpool *pgxpool.Pool) Models {
	return Models{
		Categories: CategoryModel{DBPOOL: dbpool},
	}
}

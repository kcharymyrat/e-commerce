package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Categories CategoryRepository
}

func NewRepositories(dbpool *pgxpool.Pool) Repositories {
	return Repositories{
		Categories: CategoryRepository{DBPOOL: dbpool},
	}
}

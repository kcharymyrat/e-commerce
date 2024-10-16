package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Categories   CategoryRepository
	Languages    LanguageRepository
	Translations TranslationRepository
	Users        UserRepository
}

func NewRepositories(dbpool *pgxpool.Pool) Repositories {
	return Repositories{
		Categories:   CategoryRepository{DBPOOL: dbpool},
		Languages:    LanguageRepository{DBPOOL: dbpool},
		Translations: TranslationRepository{DBPOOL: dbpool},
		Users:        UserRepository{DBPOOL: dbpool},
	}
}

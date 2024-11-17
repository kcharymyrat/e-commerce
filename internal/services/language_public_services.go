package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func GetLanguagePublicService(app *app.Application, id uuid.UUID) (*data.Language, error) {
	return app.Repositories.Languages.GetByID(id)
}

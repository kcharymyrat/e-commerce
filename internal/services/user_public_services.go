package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func GetUserByIDPublicService(app *app.Application, id uuid.UUID) (*data.User, error) {
	return GetUserByIDService(app, id)
}

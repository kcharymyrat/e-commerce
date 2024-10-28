package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateSessionService(app *app.Application, session *data.Session) error {
	return app.Repositories.Sessions.Create(session)
}

func GetSessionByRefreshTokenService(app *app.Application, refreshToken string) (*data.Session, error) {
	return app.Repositories.Sessions.GetByRefreshToken(refreshToken)
}

func GetSessionByIDService(app *app.Application, id uuid.UUID) (*data.Session, error) {
	return app.Repositories.Sessions.GetByID(id)
}

func RevokeSessionByIDService(app *app.Application, id uuid.UUID) error {
	return app.Repositories.Sessions.RevokeSessionByID(id)
}

func DeleteSessionByIDService(app *app.Application, id uuid.UUID) error {
	return app.Repositories.Sessions.DeleteByID(id)
}

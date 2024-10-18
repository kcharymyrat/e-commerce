package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateUserService(app *app.Application, user *data.User) error {
	return app.Repositories.Users.Create(user)
}

func GetUserService(app *app.Application, id uuid.UUID) (*data.User, error) {
	return app.Repositories.Users.Get(id)
}

func ListUsersService(app *app.Application, f *requests.ListUsersFilters) ([]*data.User, common.Metadata, error) {
	return app.Repositories.Users.List(f)
}

func UpdateUsersService(
	app *app.Application,
	input *requests.UpdateUserInput,
	user *data.User,
) error {
	user.Phone = input.Phone
	user.Password = input.Password
	user.FirstName = &input.FirstName
	user.LastName = &input.LastName
	user.Patronomic = &input.Patronomic
	user.Email = &input.Email
	user.IsActive = input.IsActive
	user.UpdatedByID = &input.UpdatedByID
	user.Version = input.Version

	return app.Repositories.Users.Update(user)
}

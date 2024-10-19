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

func ListUsersService(app *app.Application, f *requests.UsersAdminFilters) ([]*data.User, common.Metadata, error) {
	return app.Repositories.Users.List(f)
}

func UpdateUsersAdminService(
	app *app.Application,
	input *requests.UserAdminUpdate,
	user *data.User,
) error {
	user.Phone = input.Phone
	user.Password = input.Password
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Patronomic = input.Patronomic
	user.Email = input.Email
	user.IsActive = input.IsActive
	user.UpdatedByID = &input.UpdatedByID

	return app.Repositories.Users.Update(user)
}

func PartialUpdateUsersAdminService(
	app *app.Application,
	input *requests.UserAdminPartialUpdate,
	user *data.User,
) error {
	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.Password != nil {
		user.Password = *input.Password
	}
	if input.FirstName != nil {
		user.FirstName = input.FirstName
	}
	if input.LastName != nil {
		user.LastName = input.LastName
	}
	if input.Patronomic != nil {
		user.Patronomic = input.Patronomic
	}
	if input.Email != nil {
		user.Email = input.Email
	}
	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}
	user.UpdatedByID = &input.UpdatedByID

	return app.Repositories.Users.Update(user)
}

func DeleteUserService(app *app.Application, id uuid.UUID) error {
	return app.Repositories.Users.Delete(id)
}

func UpdateUsersSelfService(
	app *app.Application,
	input *requests.UserSelfUpdate,
	user *data.User,
) error {
	user.Phone = input.Phone
	user.Password = input.Password
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Patronomic = input.Patronomic
	user.Email = input.Email
	user.UpdatedByID = &input.UpdatedByID

	return app.Repositories.Users.Update(user)
}

func PartialUpdateUsersSelfService(
	app *app.Application,
	input *requests.UserSelfPartialUpdate,
	user *data.User,
) error {
	if input.Phone != nil {
		user.Phone = *input.Phone
	}
	if input.Password != nil {
		user.Password = *input.Password
	}
	if input.FirstName != nil {
		user.FirstName = input.FirstName
	}
	if input.LastName != nil {
		user.LastName = input.LastName
	}
	if input.Patronomic != nil {
		user.Patronomic = input.Patronomic
	}
	if input.Email != nil {
		user.Email = input.Email
	}
	user.UpdatedByID = &input.UpdatedByID

	return app.Repositories.Users.Update(user)
}

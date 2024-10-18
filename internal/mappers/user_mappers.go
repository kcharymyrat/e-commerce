package mappers

import (
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func UserCreateAdminToUser(input *requests.UserAdminCreate) *data.User {
	return &data.User{
		Phone:       input.Phone,
		Password:    input.Password,
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Patronomic:  input.Patronomic,
		Email:       input.Email,
		CreatedByID: input.CreatedByID,
	}
}

func UserToUserAdminResponse(user *data.User) *responses.UserAdminResponse {
	var res responses.UserAdminResponse
	res.ID = user.ID
	res.Phone = user.Phone
	res.FirstName = user.FirstName
	res.LastName = user.LastName
	res.Patronomic = user.Patronomic
	res.DOB = user.DOB
	res.Email = user.Email
	res.IsActive = user.IsActive
	res.IsBanned = user.IsBanned
	res.IsTrusted = user.IsTrusted
	res.InvitedByID = user.InvitedByID
	res.InvRefID = user.InvRefID
	res.InvProdRefID = user.InvProdRefID
	res.RefSignups = user.RefSignups
	res.ProdRefSignups = user.ProdRefSignups
	res.ProdRefBought = user.ProdRefBought
	res.TotalRefferals = user.TotalRefferals
	res.WholeDynDiscPercent = user.WholeDynDiscPercent
	res.DynDiscPercent = user.DynDiscPercent
	res.BonusPoints = user.BonusPoints
	res.IsStaff = user.IsStaff
	res.IsAdmin = user.IsAdmin
	res.IsSuperuser = user.IsSuperuser
	res.CreatedAt = user.CreatedAt
	res.UpdatedAt = user.UpdatedAt
	res.CreatedByID = user.CreatedByID
	res.UpdatedByID = user.UpdatedByID
	res.Version = user.Version
	return &res
}

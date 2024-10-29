package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/auth"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func CreateUserAdminHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		// TODO: authentication

		var input requests.UserAdminCreate
		err := common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		user := mappers.UserCreateAdminToUser(&input)
		err = app.Validator.Struct(user)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.CreateUserService(app, user)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				err = common.TransformPgErrToCustomError(pgErr)
				HandlePGErrors(app.Logger, localizer, w, r, err)
				return
			}
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("api/v1/users/%s", user.ID))

		res := mappers.UserToUserAdminResponse(user)
		err = common.WriteJson(w, http.StatusCreated, types.Envelope{"user": res}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func GetUsersAdminHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		user, err := services.GetUserByIDService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		res := mappers.UserToUserAdminResponse(user)
		err = common.WriteJson(w, http.StatusOK, types.Envelope{"user": res}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func ListUsersAdminHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		filters := requests.UsersAdminFilters{}
		readUserAdminQueryParams(&filters, r.URL.Query())

		err := app.Validator.Struct(filters)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			transErrs := make(map[string]string)
			for _, e := range errs {
				transErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, transErrs)
			return
		}

		users, metadata, err := services.ListUsersService(app, &filters)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		resSlice := []*responses.UserAdminResponse{}
		for _, user := range users {
			res := mappers.UserToUserAdminResponse(user)
			resSlice = append(resSlice, res)
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{
			"metadata": metadata,
			"result":   resSlice,
		}, nil)

		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}

	}
}

func UpdateUserAdminHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		input := requests.UserAdminUpdate{}
		err = common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		user, err := services.GetUserByIDService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = services.UpdateUsersAdminService(app, &input, user)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				e := common.TransformPgErrToCustomError(pgErr)
				HandlePGErrors(app.Logger, localizer, w, r, e)
				return
			}
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		res := mappers.UserToUserAdminResponse(user)
		err = common.WriteJson(w, http.StatusOK, types.Envelope{"user": res}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}

	}
}

func PartialUpdateUserAdminHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		input := requests.UserAdminPartialUpdate{}
		err = common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		user, err := services.GetUserByIDService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = services.PartialUpdateUsersAdminService(app, &input, user)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				e := common.TransformPgErrToCustomError(pgErr)
				HandlePGErrors(app.Logger, localizer, w, r, e)
				return
			}
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		res := mappers.UserToUserAdminResponse(user)
		err = common.WriteJson(w, http.StatusOK, types.Envelope{"user": res}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}

	}
}

func DeleteUserAdminHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = services.DeleteUserService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{"message": "user successfully deleted"}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func LoginUserHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		input := requests.AdminUserLoginReq{}
		err := common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		user, err := services.GetUserByPhoneService(app, input.Phone)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		ok, err := auth.IsPasswordInputMatching(input.Password, user.PasswordHash)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}
		if !ok {
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		if !user.IsActive || user.IsBanned {
			app.Logger.Info().Msg("user is not active or banned")
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		accessToken, accessClaims, err := auth.GenerateJWT(
			user.ID,
			user.Phone,
			user.FirstName,
			user.LastName,
			user.Patronomic,
			user.IsActive,
			user.IsBanned,
			user.IsStaff,
			user.IsAdmin,
			user.IsSuperuser,
			5*time.Minute,
			app.Config.SecretKey,
			app.Logger,
		)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		refreshToken, refreshClaims, err := auth.GenerateJWT(
			user.ID,
			user.Phone,
			user.FirstName,
			user.LastName,
			user.Patronomic,
			user.IsActive,
			user.IsBanned,
			user.IsStaff,
			user.IsAdmin,
			user.IsSuperuser,
			48*time.Hour,
			app.Config.SecretKey,
			app.Logger,
		)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		session := data.Session{
			ID:           uuid.MustParse(refreshClaims.RegisteredClaims.ID),
			UserPhone:    user.Phone,
			RefreshToken: refreshToken,
			IsRevoked:    false,
			ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
		}
		err = services.CreateSessionService(app, &session)
		if err != nil {
			app.Logger.Error().Err(err).Msg("failed to create session")
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		res := responses.LoginResponse{
			SessionID:             session.ID,
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
			RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
			User: responses.ShortUserResponse{
				ID:          user.ID,
				Phone:       user.Phone,
				FirstName:   user.FirstName,
				LastName:    user.LastName,
				Patronomic:  user.Patronomic,
				IsActive:    user.IsActive,
				IsBanned:    user.IsBanned,
				IsStaff:     user.IsStaff,
				IsAdmin:     user.IsAdmin,
				IsSuperuser: user.IsSuperuser,
			},
		}
		err = common.WriteJson(w, http.StatusOK, types.Envelope{"result": res}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}

	}
}

func LogoutUserHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		authorization := r.Header.Get("Authorization")
		access_token := strings.TrimPrefix(authorization, "Bearer ")

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		session, err := services.GetSessionByIDService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		accessClaims, err := auth.ParseJWT(access_token, app.Config.SecretKey, app.Logger)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		if accessClaims.Phone != session.UserPhone {
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		err = services.DeleteSessionByIDService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{"message": "session successfully deleted"}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

// TODO: revoke token when user is banned
func RenewAccessTokenReqHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		authorization := r.Header.Get("Authorization")
		access_token := strings.TrimPrefix(authorization, "Bearer ")

		input := requests.RenewAccessTokenReq{}
		err := common.ReadJSON(w, r, input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		refreshClaims, err := auth.ParseJWT(input.RefreshToken, app.Config.SecretKey, app.Logger)
		if err != nil {
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		accessClaims, err := auth.ParseJWT(access_token, app.Config.SecretKey, app.Logger)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		if accessClaims.Phone != refreshClaims.Phone {
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		session, err := services.GetSessionByRefreshTokenService(app, input.RefreshToken)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				app.Logger.Error().Err(err).Msg("session not found")
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				app.Logger.Error().Err(err).Msg("failed to get session")
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		if session.IsRevoked {
			app.Logger.Error().Err(err).Msg("session is revoked")
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		if session.ExpiresAt.Before(time.Now()) {
			app.Logger.Error().Err(err).Msg("session expired")
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		if session.UserPhone != refreshClaims.Phone {
			app.Logger.Error().Err(err).Msg("session phone mismatch")
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		accessToken, accessClaims, err := auth.GenerateJWT(
			refreshClaims.UserID,
			refreshClaims.Phone,
			refreshClaims.FirstName,
			refreshClaims.LastName,
			refreshClaims.Patronomic,
			refreshClaims.IsActive,
			refreshClaims.IsBanned,
			refreshClaims.IsStaff,
			refreshClaims.IsAdmin,
			refreshClaims.IsSuperuser,
			5*time.Minute,
			app.Config.SecretKey,
			app.Logger,
		)
		if err != nil {
			app.Logger.Error().Err(err).Msg("failed to generate access token")
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		res := responses.RenewAccessTokenResponse{
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,
		}
		err = common.WriteJson(w, http.StatusOK, types.Envelope{"result": res}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func RevokeSessionByIDHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		authorization := r.Header.Get("Authorization")
		access_token := strings.TrimPrefix(authorization, "Bearer ")

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		session, err := services.GetSessionByIDService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.UnauthorizedResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		accessClaims, err := auth.ParseJWT(access_token, app.Config.SecretKey, app.Logger)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		if accessClaims.Phone != session.UserPhone {
			common.UnauthorizedResponse(app.Logger, localizer, w, r)
			return
		}

		err = services.RevokeSessionByIDService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusNoContent, nil, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

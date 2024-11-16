package handlers

import (
	"errors"
	"net/http"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/auth"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func LoginWithPasswordPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		input := &requests.UserLoginReq{}
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
				common.BadRequestResponse(app.Logger, localizer, w, r, err)
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

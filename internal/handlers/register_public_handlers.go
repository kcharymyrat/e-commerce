package handlers

import (
	"fmt"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func RegisterUserWithPasswordPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		var input requests.UserPasswordRegisterReq
		err := common.ReadJSON(w, r, &input)
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

		user := data.User{
			Phone:    input.Phone,
			Password: input.Password,
			IsActive: true,
		}
		err = services.CreateUserService(app, &user)

		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				err = common.TransformPgErrToCustomError(pgErr)
				HandlePGErrors(app.Logger, localizer, w, r, err)
				return
			}
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		// TODO: Handle referrals in here

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("api/v1/users/%s", user.ID))

		res := mappers.UserToUserSelfResponse(&user)
		err = common.WriteJson(w, http.StatusCreated, types.Envelope{"user": res}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}

	}
}

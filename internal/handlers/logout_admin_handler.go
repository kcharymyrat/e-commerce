package handlers

import (
	"errors"
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/auth"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func LogoutAdminHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

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

		accessClaims := r.Context().Value(types.UserClaimsKey{}).(*auth.UserClaims)

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

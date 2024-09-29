package handlers

import (
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func HealthcheckHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		data := map[string]string{
			"status":      "available",
			"environment": app.Config.Env,
		}

		err := common.WriteJson(w, http.StatusOK, common.Envelope{"data": data}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

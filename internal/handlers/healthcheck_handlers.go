package handlers

import (
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// @Summary Get health check
// @Tags Healthcheck
// @Description Returns the status of the API
// @ID healthcheck
// @Accept json
// @Produce json
// @Router /api/v1/healthcheck [get]
// @Success 200 {object} types.Envelope
// @Failure 500 {object} types.ErrorResponse
func HealthcheckHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		data := map[string]string{
			"status":      "available",
			"environment": app.Config.Env,
		}

		err := common.WriteJson(w, http.StatusOK, types.Envelope{"data": data}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

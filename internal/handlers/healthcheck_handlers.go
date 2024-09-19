package handlers

import (
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
)

func HealthcheckHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"status":      "available",
			"environment": app.Config.Env,
		}

		err := common.WriteJson(w, http.StatusOK, common.Envelope{"data": data}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

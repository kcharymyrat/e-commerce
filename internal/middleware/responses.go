package middleware

import (
	"fmt"
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/rs/zerolog"
)

func NotFound(logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := "the requested resource could not be found"
		envel := common.Envelope{"error": message}

		err := common.WriteJson(w, http.StatusNotFound, envel, nil)
		if err != nil {
			logger.Error().Err(err).Str("url", r.URL.String()).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func MethodNotAllowed(logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
		envel := common.Envelope{"error": message}

		err := common.WriteJson(w, http.StatusMethodNotAllowed, envel, nil)
		if err != nil {
			logger.Error().Err(err).Str("url", r.URL.String()).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

package middleware

import (
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
)

func NotFound(logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)
		message, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "not_found",
		})

		if err != nil {
			logger.Error().Err(err).Str("url", r.URL.String()).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
		}

		envel := common.Envelope{"error": message}

		err = common.WriteJson(w, http.StatusNotFound, envel, nil)
		if err != nil {
			logger.Error().Err(err).Str("url", r.URL.String()).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func MethodNotAllowed(logger *zerolog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		message, err := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: "method_not_allowed",
			TemplateData: map[string]interface{}{
				"method": r.Method,
			},
		})

		if err != nil {
			logger.Error().Err(err).Str("url", r.URL.String()).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		envel := common.Envelope{"error": message}

		err = common.WriteJson(w, http.StatusMethodNotAllowed, envel, nil)
		if err != nil {
			logger.Error().Err(err).Str("url", r.URL.String()).Msg("error")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

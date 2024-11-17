package common

import (
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
)

func LogError(logger *zerolog.Logger, r *http.Request, err error) {
	logger.Error().Err(err).Str("url", r.URL.String()).Msg("error")
}

func ErrorResponse(
	logger *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message interface{},
) {
	envel := types.ErrorResponse{
		Code:  status,
		Error: message.(string),
	}

	err := WriteErrorJson(w, status, envel, nil)
	if err != nil {
		LogError(logger, r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	LogError(logger, r, err)

	message, e := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "server_error",
	})

	if e != nil {
		ErrorResponse(logger, w, r, http.StatusInternalServerError, e.Error())
		return
	}

	ErrorResponse(logger, w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
) {
	message, e := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "not_found",
	})

	if e != nil {
		ErrorResponse(logger, w, r, http.StatusInternalServerError, e.Error())
		return
	}

	ErrorResponse(logger, w, r, http.StatusNotFound, message)
}

func UnauthorizedResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
) {
	message, e := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "unauthorized",
	})

	if e != nil {
		ErrorResponse(logger, w, r, http.StatusInternalServerError, e.Error())
		return
	}

	ErrorResponse(logger, w, r, http.StatusUnauthorized, message)
}

func MethodNotAllowedResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
) {
	message, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "method_not_allowed",
		TemplateData: map[string]interface{}{
			"method": r.Method, // Pass the HTTP method as template data
		},
	})

	if err != nil {
		ErrorResponse(logger, w, r, http.StatusInternalServerError, err.Error())
		return
	}

	ErrorResponse(logger, w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	message, e := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "bad_request",
		TemplateData: map[string]interface{}{
			"error": err.Error(),
		},
	})

	if e != nil {
		ErrorResponse(logger, w, r, http.StatusInternalServerError, e.Error())
		return
	}

	ErrorResponse(logger, w, r, http.StatusBadRequest, message)
}

func EditConflictResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
) {
	message, e := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "edit_conflict",
	})

	if e != nil {
		ErrorResponse(logger, w, r, http.StatusInternalServerError, e.Error())
		return
	}

	ErrorResponse(logger, w, r, http.StatusConflict, message)
}

func RateLimitExceedResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
) {
	message, e := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "rate_limit_exceeded",
	})

	if e != nil {
		ErrorResponse(logger, w, r, http.StatusInternalServerError, e.Error())
		return
	}
	ErrorResponse(logger, w, r, http.StatusTooManyRequests, message)
}

func FailedValidationResponse(
	logger *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	ErrorResponse(logger, w, r, http.StatusUnprocessableEntity, errors)
}

package common

import (
	"fmt"
	"net/http"

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
	envel := envelope{"error": message}

	err := WriteJson(w, status, envel, nil)
	if err != nil {
		LogError(logger, r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(
	logger *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	LogError(logger, r, err)
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(logger, w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(
	logger *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
) {
	message := "the requested resource could not be found"
	ErrorResponse(logger, w, r, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(
	logger *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	ErrorResponse(logger, w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(
	logger *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	ErrorResponse(logger, w, r, http.StatusBadRequest, err.Error())
}

func FailedValidationResponse(
	logger *zerolog.Logger,
	w http.ResponseWriter,
	r *http.Request,
	errors map[string]string,
) {
	ErrorResponse(logger, w, r, http.StatusUnprocessableEntity, errors)
}

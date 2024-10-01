package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")
var ErrInvalidSlug = errors.New("invalid slug")

var (
	ErrIntegrityConstraintViolation = errors.New("integrity constraint violation")
	ErrRestrictViolation            = errors.New("restrict violation")
	ErrNotNullViolation             = errors.New("not null constraint violation")
	ErrForeignKeyViolation          = errors.New("foreign key violation")
	ErrUniqueViolation              = errors.New("unique constraint violation")
	ErrCheckViolation               = errors.New("check constraint violation")
	ErrExclusionViolation           = errors.New("exlusion violation")

	ErrStringDataTruncation   = errors.New("string data truncation")
	ErrNumericValueOutOfRange = errors.New("numeric value out of range")
	ErrInvalidDatetimeFormat  = errors.New("invalid datatime format")
)

func TransformPgErrToCustomError(pgErr *pgconn.PgError) error {
	switch pgErr.Code {
	case constants.IntegrityConstraintViolation:
		return fmt.Errorf("%w: %s", ErrIntegrityConstraintViolation, pgErr.Detail)
	case constants.RestrictViolation:
		return fmt.Errorf("%w: %s", ErrRestrictViolation, pgErr.Detail)
	case constants.NotNullViolation:
		return fmt.Errorf("%w: %s", ErrNotNullViolation, pgErr.Detail)
	case constants.ForeignKeyViolation:
		return fmt.Errorf("%w: %s", ErrForeignKeyViolation, pgErr.Detail)
	case constants.UniqueViolation:
		return fmt.Errorf("%w: %s", ErrUniqueViolation, pgErr.Detail)
	case constants.ExclusionViolation:
		return fmt.Errorf("%w: %s", ErrCheckViolation, pgErr.Detail)
	case constants.CheckViolation:
		return fmt.Errorf("%w: %s", ErrCheckViolation, pgErr.Detail)
	case constants.StringDataRightTruncationDataException:
		return fmt.Errorf("%w: %s", ErrStringDataTruncation, pgErr.Detail)
	case constants.NumericValueOutOfRange:
		return fmt.Errorf("%w: %s", ErrNumericValueOutOfRange, pgErr.Detail)
	case constants.InvalidDatetimeFormat:
		return fmt.Errorf("%w: %s", ErrInvalidDatetimeFormat, pgErr.Detail)
	default:
		return fmt.Errorf("%w: %s", pgErr, pgErr.Detail)
	}
}

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
	envel := types.Envelope{"error": message}

	err := WriteJson(w, status, envel, nil)
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

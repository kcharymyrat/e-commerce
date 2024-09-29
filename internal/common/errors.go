package common

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")

var (
	ErrIntegrityConstraintViolation = fmt.Errorf("%s :integrity constraint violation", IntegrityConstraintViolation)
	ErrRestrictViolation            = fmt.Errorf("%s :restrict violation", RestrictViolation)
	ErrNotNullViolation             = fmt.Errorf("%s :not null constraint violation", NotNullViolation)
	ErrForeignKeyViolation          = fmt.Errorf("%s :foreign key violation", ForeignKeyViolation)
	ErrUniqueViolation              = fmt.Errorf("%s :unique constraint violation", UniqueViolation)
	ErrCheckViolation               = fmt.Errorf("%s :check constraint violation", CheckViolation)
	ErrExclusionViolation           = fmt.Errorf("%s :exlusion violation", ExclusionViolation)

	ErrStringDataTruncation   = fmt.Errorf("%s :string data truncation", StringDataRightTruncationDataException)
	ErrNumericValueOutOfRange = fmt.Errorf("%s :numeric value out of range", NumericValueOutOfRange)
	ErrInvalidDatetimeFormat  = fmt.Errorf("%s :invalid datatime format", InvalidDatetimeFormat)
)

func TransformPgErrToCustomError(pgErr *pgconn.PgError) error {
	switch pgErr.Code {
	case IntegrityConstraintViolation:
		return fmt.Errorf("%w: %s", ErrIntegrityConstraintViolation, pgErr.Detail)
	case RestrictViolation:
		return fmt.Errorf("%w: %s", ErrRestrictViolation, pgErr.Detail)
	case NotNullViolation:
		return fmt.Errorf("%w: %s", ErrNotNullViolation, pgErr.Detail)
	case ForeignKeyViolation:
		return fmt.Errorf("%w: %s", ErrForeignKeyViolation, pgErr.Detail)
	case UniqueViolation:
		return fmt.Errorf("%w: %s", ErrUniqueViolation, pgErr.Detail)
	case ExclusionViolation:
		return fmt.Errorf("%w: %s", ErrCheckViolation, pgErr.Detail)
	case CheckViolation:
		return fmt.Errorf("%w: %s", ErrCheckViolation, pgErr.Detail)
	case StringDataRightTruncationDataException:
		return fmt.Errorf("%w: %s", ErrStringDataTruncation, pgErr.Detail)
	case NumericValueOutOfRange:
		return fmt.Errorf("%w: %s", ErrNumericValueOutOfRange, pgErr.Detail)
	case InvalidDatetimeFormat:
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
	envel := Envelope{"error": message}

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

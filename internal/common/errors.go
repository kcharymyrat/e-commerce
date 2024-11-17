package common

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/internal/constants"
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

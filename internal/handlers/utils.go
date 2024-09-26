package handlers

import (
	"errors"
	"net/http"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
)

func HandleCategoryServiceErrors(w http.ResponseWriter, r *http.Request, app *app.Application, err error) {
	switch {
	case errors.Is(err, common.ErrIntegrityConstraintViolation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrIntegrityConstraintViolation)
	case errors.Is(err, common.ErrRestrictViolation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrRestrictViolation)
	case errors.Is(err, common.ErrNotNullViolation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrNotNullViolation)
	case errors.Is(err, common.ErrForeignKeyViolation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrForeignKeyViolation)
	case errors.Is(err, common.ErrUniqueViolation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrUniqueViolation)
	case errors.Is(err, common.ErrCheckViolation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrCheckViolation)
	case errors.Is(err, common.ErrExclusionViolation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrExclusionViolation)
	case errors.Is(err, common.ErrStringDataTruncation):
		common.BadRequestResponse(app.Logger, w, r, common.ErrStringDataTruncation)
	case errors.Is(err, common.ErrNumericValueOutOfRange):
		common.BadRequestResponse(app.Logger, w, r, common.ErrNumericValueOutOfRange)
	case errors.Is(err, common.ErrInvalidDatetimeFormat):
		common.BadRequestResponse(app.Logger, w, r, common.ErrInvalidDatetimeFormat)
	default:
		common.ServerErrorResponse(app.Logger, w, r, err)
	}
}

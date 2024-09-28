package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/kcharymyrat/e-commerce/api/requests"
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

func readCategoryQueryParameters(input *requests.ListCategoriesInput, qs url.Values) {
	input.Names = common.ReadQueryCSStrs(qs, "names")
	input.Slugs = common.ReadQueryCSStrs(qs, "slugs")
	input.ParentIDs = common.ReadQueryCSUUIDs(qs, "parent_ids")
	input.Search = common.ReadQueryStr(qs, "search")
	input.CreatedAtFrom = common.ReadQueryTime(qs, "created_at_from")
	input.CreatedAtUpTo = common.ReadQueryTime(qs, "created_at_up_to")
	input.UpdatedAtFrom = common.ReadQueryTime(qs, "updated_at_from")
	input.UpdatedAtUpTo = common.ReadQueryTime(qs, "updated_at_up_to")
	input.CreatedByIDs = common.ReadQueryCSUUIDs(qs, "created_by_ids")
	input.UpdatedByIDs = common.ReadQueryCSUUIDs(qs, "updated_by_ids")
	input.Sorts = common.ReadQueryCSStrs(qs, "sorts")
	input.SortSafeList = []string{
		"id", "name", "created_at", "updated_at", "-id", "-name", "-created_at", "-updated_at",
	}
	input.Page = common.ReadQueryInt(qs, "page")
	input.PageSize = common.ReadQueryInt(qs, "page_size")
}

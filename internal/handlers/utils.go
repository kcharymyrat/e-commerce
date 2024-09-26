package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/filters"
	"github.com/kcharymyrat/e-commerce/internal/validator"
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

func readAndValidateQueryParameters(input *requests.ListCategoriesInput, qs url.Values, v *validator.Validator) {
	input.Names = common.ReadQueryCSStrs(qs, "names")
	input.Slugs = common.ReadQueryCSStrs(qs, "slugs")
	input.ParentIDs = common.ReadQueryCSUUIDs(qs, "parent_ids", v)
	input.Search = common.ReadQueryStr(qs, "search")
	input.CreatedAtFrom = common.ReadQueryTime(qs, "created_at_from", v)
	input.CreatedAtUpTo = common.ReadQueryTime(qs, "created_at_up_to", v)
	input.UpdatedAtFrom = common.ReadQueryTime(qs, "updated_at_from", v)
	input.UpdatedAtUpTo = common.ReadQueryTime(qs, "updated_at_up_to", v)
	input.CreatedByIDs = common.ReadQueryCSUUIDs(qs, "created_by_ids", v)
	input.UpdatedByIDs = common.ReadQueryCSUUIDs(qs, "updated_by_ids", v)
	input.Sorts = common.ReadQueryCSStrs(qs, "sorts")
	input.SortSafeList = []string{
		"id", "name", "created_at", "updated_at", "-id", "-name", "-created_at", "-updated_at",
	}
	input.Page = common.ReadQueryInt(qs, "page", v)
	input.PageSize = common.ReadQueryInt(qs, "page_size", v)
}

func filtersValidation(input *requests.ListCategoriesInput, v *validator.Validator) {
	filters.ValidateSearchFilters(v, input.SearchFilters)
	filters.ValidateCreatedUpdatedAtFilters(v, input.CreatedUpdatedAtFilters)
	filters.ValidateCreatedUpdatedByFilters(v, input.CreatedUpdatedByFilters)
	filters.ValidateSortFilters(v, input.SortListFilters)
	filters.ValidatePaginationFilters(v, input.PaginationFilters)
}

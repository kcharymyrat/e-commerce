package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/rs/zerolog"
)

func serviceBadRequestResponse(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
	messageId string,
	err error,
) {
	message, e := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageId,
		TemplateData: map[string]interface{}{
			"details": err.Error(),
		},
	})

	if e != nil {
		common.ErrorResponse(logger, w, r, http.StatusInternalServerError, e.Error())
		return
	}

	common.ErrorResponse(logger, w, r, http.StatusBadRequest, message)
}

func HandlePGErrors(
	logger *zerolog.Logger,
	localizer *i18n.Localizer,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	switch {
	case errors.Is(err, common.ErrIntegrityConstraintViolation):
		messageId := "integrity_constraint_violation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrIntegrityConstraintViolation)
	case errors.Is(err, common.ErrRestrictViolation):
		messageId := "restrict_violation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrRestrictViolation)
	case errors.Is(err, common.ErrNotNullViolation):
		messageId := "restrict_violation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrNotNullViolation)
	case errors.Is(err, common.ErrForeignKeyViolation):
		messageId := "foreign_key_violation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrForeignKeyViolation)
	case errors.Is(err, common.ErrUniqueViolation):
		messageId := "unique_violation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrUniqueViolation)
	case errors.Is(err, common.ErrCheckViolation):
		messageId := "check_violation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrCheckViolation)
	case errors.Is(err, common.ErrExclusionViolation):
		messageId := "exclusion_violation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrExclusionViolation)
	case errors.Is(err, common.ErrStringDataTruncation):
		messageId := "string_data_truncation"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrStringDataTruncation)
	case errors.Is(err, common.ErrNumericValueOutOfRange):
		messageId := "numeric_value_out_of_range"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrNumericValueOutOfRange)
	case errors.Is(err, common.ErrInvalidDatetimeFormat):
		messageId := "invalid_datetime_format"
		serviceBadRequestResponse(logger, localizer, w, r, messageId, common.ErrInvalidDatetimeFormat)
	default:
		common.ServerErrorResponse(logger, localizer, w, r, err)
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

func readLanguageQueryParameters(input *requests.ListLanguagesInput, qs url.Values) {
	input.Page = common.ReadQueryInt(qs, "page")
	input.PageSize = common.ReadQueryInt(qs, "page_size")
}

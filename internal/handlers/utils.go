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

func readCategoryAdminQueryParams(input *requests.CategoriesAdminFilters, qs url.Values) {
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

func readLanguageAdminQueryParams(input *requests.LanguagesAdminFilters, qs url.Values) {
	input.Page = common.ReadQueryInt(qs, "page")
	input.PageSize = common.ReadQueryInt(qs, "page_size")
}

func readTranslationAdminQueryParams(input *requests.TranslationsAdminFilters, qs url.Values) {
	input.LanguageCodes = common.ReadQueryCSStrs(qs, "language_codes")
	input.TableNames = common.ReadQueryCSStrs(qs, "table_names")
	input.FieldNames = common.ReadQueryCSStrs(qs, "field_names")
	input.EntityIDs = common.ReadQueryCSUUIDs(qs, "entity_ids")
	input.Search = common.ReadQueryStr(qs, "search")
	input.CreatedAtFrom = common.ReadQueryTime(qs, "created_at_from")
	input.CreatedAtUpTo = common.ReadQueryTime(qs, "created_at_up_to")
	input.UpdatedAtFrom = common.ReadQueryTime(qs, "updated_at_from")
	input.UpdatedAtUpTo = common.ReadQueryTime(qs, "updated_at_up_to")
	input.CreatedByIDs = common.ReadQueryCSUUIDs(qs, "created_by_ids")
	input.UpdatedByIDs = common.ReadQueryCSUUIDs(qs, "updated_by_ids")
	input.Sorts = common.ReadQueryCSStrs(qs, "sorts")
	input.SortSafeList = []string{
		"id", "language_code", "entity_id", "table_name", "-id", "-language_code", "-entity_ids", "-table_name",
	}
	input.Page = common.ReadQueryInt(qs, "page")
	input.PageSize = common.ReadQueryInt(qs, "page_size")
}

func readUserAdminQueryParams(input *requests.UsersAdminFilters, qs url.Values) {
	input.ID = common.ReadQueryUUID(qs, "id")
	input.Phone = common.ReadQueryStr(qs, "phone")
	input.Email = common.ReadQueryStr(qs, "email")
	input.IsActice = common.ReadQueryBool(qs, "is_active")
	input.IsBanned = common.ReadQueryBool(qs, "is_banned")
	input.IsTrusted = common.ReadQueryBool(qs, "is_trusted")
	input.IsInvited = common.ReadQueryBool(qs, "is_invited")
	input.RefSignupsFrom = common.ReadQueryInt(qs, "ref_signups_from")
	input.RefSignupsTo = common.ReadQueryInt(qs, "ref_signups_to")
	input.ProdRefBoughtFrom = common.ReadQueryInt(qs, "prod_ref_bought_from")
	input.ProdRefBoughtTo = common.ReadQueryInt(qs, "prod_ref_bought_to")
	input.ProdRefBoughtFrom = common.ReadQueryInt(qs, "prod_ref_bought_from")
	input.ProdRefBoughtTo = common.ReadQueryInt(qs, "prod_ref_bought_to")
	input.WholeDynDiscPercentFrom = common.ReadQueryDecimal(qs, "whole_ddp_from")
	input.WholeDynDiscPercentTo = common.ReadQueryDecimal(qs, "whole_ddp_to")
	input.DynDiscPercentFrom = common.ReadQueryDecimal(qs, "ddp_from")
	input.DynDiscPercentTo = common.ReadQueryDecimal(qs, "ddp_to")
	input.BonusPointsFrom = common.ReadQueryDecimal(qs, "bonus_from")
	input.BonusPointsTo = common.ReadQueryDecimal(qs, "bonus_to")
	input.IsStaff = common.ReadQueryBool(qs, "is_staff")
	input.IsAdmin = common.ReadQueryBool(qs, "is_admin")
	input.IsSuperuser = common.ReadQueryBool(qs, "is_superuser")
	input.Search = common.ReadQueryStr(qs, "search")
	input.CreatedAtFrom = common.ReadQueryTime(qs, "created_at_from")
	input.CreatedAtUpTo = common.ReadQueryTime(qs, "created_at_up_to")
	input.UpdatedAtFrom = common.ReadQueryTime(qs, "updated_at_from")
	input.UpdatedAtUpTo = common.ReadQueryTime(qs, "updated_at_up_to")
	input.CreatedByIDs = common.ReadQueryCSUUIDs(qs, "created_by_ids")
	input.UpdatedByIDs = common.ReadQueryCSUUIDs(qs, "updated_by_ids")
	input.Sorts = common.ReadQueryCSStrs(qs, "sorts")
	input.SortSafeList = []string{
		"id", "phone", "email", "is_active", "is_banned", "is_trusted", "is_invited",
		"ref_signups", "prod_ref_bought", "whole_ddp", "ddp", "bonus", "is_staff",
		"created_at", "updated_at",
		"-id", "-phone", "-email", "-is_active", "-is_banned", "-is_trusted", "-is_invited",
		"-ref_signups", "-prod_ref_bought", "-whole_ddp", "-ddp", "-bonus", "-is_staff",
		"-created_at", "-updated_at",
	}
	input.Page = common.ReadQueryInt(qs, "page")
	input.PageSize = common.ReadQueryInt(qs, "page_size")
}

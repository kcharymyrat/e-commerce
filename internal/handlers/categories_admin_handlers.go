package handlers

import (
	"errors"
	"fmt"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/constants"
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/types"
	"github.com/kcharymyrat/e-commerce/internal/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func CreateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		var categoryInput requests.CategoryAdminCreate
		err := common.ReadJSON(w, r, &categoryInput)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		category := mappers.CreateCategoryInputToCategoryMapper(&categoryInput)
		err = app.Validator.Struct(category)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.CreateCategoryService(app, category)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				err = common.TransformPgErrToCustomError(pgErr)
				HandlePGErrors(app.Logger, localizer, w, r, err)
				return
			}
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/api/v1/%v", category.ID))

		categoryResponse := mappers.CategoryToCategoryManagerResponseMapper(category)

		err = common.WriteJson(w, http.StatusCreated, types.Envelope{"category": categoryResponse}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func GetCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang_code := common.GetAcceptLanguageHeader(r)

		// valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		slug, err := common.ReadSlugParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		category, err := services.GetCategoryBySlugService(app, slug)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		categoryManagerResponse := mappers.CategoryToCategoryManagerResponseMapper(category)

		trMapWrapper := make(map[string]map[string]string)

		nameFieldTrMap := make(map[string]string)
		name_field_tr, name_value_tr, err := utils.GetTranslatedFieldNameAndValue(app, category.ID, lang_code, "name")
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}
		nameFieldTrMap["field_name"] = name_field_tr
		nameFieldTrMap["field_value"] = name_value_tr
		trMapWrapper["name"] = nameFieldTrMap

		descFieldTrMap := make(map[string]string)
		desc_field_tr, desc_value_tr, err := utils.GetTranslatedFieldNameAndValue(app, category.ID, lang_code, "description")
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}
		descFieldTrMap["field_name"] = desc_field_tr
		descFieldTrMap["field_value"] = desc_value_tr
		trMapWrapper["description"] = nameFieldTrMap

		err = common.WriteJson(w, http.StatusOK, types.Envelope{
			"category":     categoryManagerResponse,
			"translations": trMapWrapper,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func ListCategoriesManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang_code := common.GetAcceptLanguageHeader(r)

		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		filters := requests.CategoriesAdminFilters{}

		readCategoryAdminQueryParams(&filters, r.URL.Query())

		err := app.Validator.Struct(filters)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		categories, metadata, err := services.ListCategoriesService(app, &filters)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		catWithTransResponses := make([]*responses.CategoryWithTranslationsAdminResponse, len(categories))
		for _, category := range categories {
			cat := mappers.CategoryToCategoryManagerResponseMapper(category)

			trans := make(map[string]map[string]string)
			nameFieldTrMap := make(map[string]string)
			name_field_tr, name_value_tr, err := utils.GetTranslatedFieldNameAndValue(app, category.ID, lang_code, "name")
			if err != nil {
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}
			nameFieldTrMap["field_name"] = name_field_tr
			nameFieldTrMap["field_value"] = name_value_tr
			trans["name"] = nameFieldTrMap

			descFieldTrMap := make(map[string]string)
			desc_field_tr, desc_value_tr, err := utils.GetTranslatedFieldNameAndValue(app, category.ID, lang_code, "description")
			if err != nil {
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}
			descFieldTrMap["field_name"] = desc_field_tr
			descFieldTrMap["field_value"] = desc_value_tr
			trans["description"] = nameFieldTrMap

			catWithTrans := responses.CategoryWithTranslationsAdminResponse{
				Category:     *cat,
				Translations: trans,
			}

			catWithTransResponses = append(catWithTransResponses, &catWithTrans)
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{
			"metadata": metadata,
			"results":  catWithTransResponses,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func UpdateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		slug, err := common.ReadSlugParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		category, err := services.GetCategoryBySlugService(app, slug)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrEditConflict):
				common.EditConflictResponse(app.Logger, localizer, w, r)
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		input := requests.CategoryAdminUpdate{}

		err = common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.UpdateCategoryService(app, &input, category)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				err = common.TransformPgErrToCustomError(pgErr)
				HandlePGErrors(app.Logger, localizer, w, r, err)
				return
			}
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{"category": category}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func PartialUpdateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		slug, err := common.ReadSlugParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		category, err := services.GetCategoryBySlugService(app, slug)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrEditConflict):
				common.EditConflictResponse(app.Logger, localizer, w, r)
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		input := requests.CategoryAdminPartialUpdate{}

		err = common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = app.Validator.Struct(category)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		err = services.PartialUpdateCategoryService(app, &input, category)
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok {
				err = common.TransformPgErrToCustomError(pgErr)
				HandlePGErrors(app.Logger, localizer, w, r, err)
				return
			}
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, types.Envelope{"category": category}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func DeleteCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		slug, err := common.ReadSlugParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		err = services.DeleteCategoryServiceBySlug(app, slug)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		// TODO: Needs localiztions
		err = common.WriteJson(w, http.StatusOK, types.Envelope{"message": "category successfully deleted"}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

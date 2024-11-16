package handlers

import (
	"errors"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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

func GetCategoryPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang_code := common.GetAcceptLanguageHeader(r)

		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
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

		categoryPublicResponse := mappers.CategoryToCategoryPublicResponseMapper(category)

		trMapWrapper := make(map[string]map[string]string)

		nameFieldTrMap := make(map[string]string)
		name_field_tr, name_value_tr, err := utils.GetTranslationMap(app, category.ID, lang_code, "name")
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}
		nameFieldTrMap["field_name"] = name_field_tr
		nameFieldTrMap["field_value"] = name_value_tr
		trMapWrapper["name"] = nameFieldTrMap

		descFieldTrMap := make(map[string]string)
		desc_field_tr, desc_value_tr, err := utils.GetTranslationMap(app, category.ID, lang_code, "description")
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}
		descFieldTrMap["field_name"] = desc_field_tr
		descFieldTrMap["field_value"] = desc_value_tr
		trMapWrapper["description"] = nameFieldTrMap

		err = common.WriteJson(w, http.StatusOK, types.Envelope{
			"category":     categoryPublicResponse,
			"translations": trMapWrapper,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

// @Summary List all categories
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {object} types.Envelope
// @Failure 500 {object} types.Envelope
// @Failure 422 {object} types.Envelope
// @Param filters query requests.CategoriesAdminFilters true "Filters"
// @Router /api/v1/categories [get]
func ListCategoriesPublicHandler(app *app.Application) http.HandlerFunc {
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

		catWithTransResponses := make([]*responses.CategoryWithTranslationsPublicResponse, len(categories))
		for _, category := range categories {
			cat := mappers.CategoryToCategoryPublicResponseMapper(category)

			trans := make(map[string]map[string]string)
			nameFieldTrMap := make(map[string]string)
			name_field_tr, name_value_tr, err := utils.GetTranslationMap(app, category.ID, lang_code, "name")
			if err != nil {
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}
			nameFieldTrMap["field_name"] = name_field_tr
			nameFieldTrMap["field_value"] = name_value_tr
			trans["name"] = nameFieldTrMap

			descFieldTrMap := make(map[string]string)
			desc_field_tr, desc_value_tr, err := utils.GetTranslationMap(app, category.ID, lang_code, "description")
			if err != nil {
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}
			descFieldTrMap["field_name"] = desc_field_tr
			descFieldTrMap["field_value"] = desc_value_tr
			trans["description"] = nameFieldTrMap

			catWithTrans := responses.CategoryWithTranslationsPublicResponse{
				Category:     *cat,
				Translations: trans,
			}

			catWithTransResponses = append(catWithTransResponses, &catWithTrans)
		}

		// Write the response as JSON
		err = common.WriteJson(w, http.StatusOK, types.Envelope{
			"metadata": metadata,
			"results":  catWithTransResponses,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

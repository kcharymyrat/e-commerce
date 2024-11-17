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
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// @Summary List categories
// @Description List categories with pagination and filters
// @Tags categories
// @Param Accept-Language header string false "Languages: en, ru, tk"
// @Param filters query requests.CategoriesAdminFilters true "Filters"
// @Produce json
// @Router /api/v1/categories [get]
// @Success 200 {object} types.PaginatedResponse[responses.CategoryPublicResponse]
// @Failure 500 {object} types.ErrorResponse
// @Failure 422 {object} types.ErrorResponse
func ListCategoriesPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		langCode := common.GetAcceptLanguageHeader(r)

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

		catsWithTrs, metadata, err := services.ListCategoriesPublicService(app, &filters, langCode)
		catWithTransResponses := make([]*types.DetailResponse[responses.CategoryPublicResponse], len(catsWithTrs))
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		for _, catWithTrs := range catsWithTrs {
			categoryPublicResponse := mappers.CategoryToCategoryPublicResponseMapper(catWithTrs.Category)
			detailResponse := types.NewDetailResponse(categoryPublicResponse, catWithTrs.Translations)
			catWithTransResponses = append(catWithTransResponses, detailResponse)
		}

		// Write the response as JSON
		paginatedRes := types.PaginatedResponse[responses.CategoryPublicResponse]{
			Metadata: metadata,
			Results:  catWithTransResponses,
		}
		err = common.WritePaginatedJson(w, http.StatusOK, paginatedRes, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

// @Summary Get category by slug
// @Description Get specific category details by slug
// @Tags categories
// @Param Accept-Language header string false "Languages: en, ru, tk"
// @Param slug path string true "Category Slug"
// @Accept multipart/form-data
// @Produce json
// @Router /api/v1/categories/{slug} [get]
// @Success 200 {object} types.DetailResponse[responses.CategoryPublicResponse]
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
func GetCategoryPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		langCode := common.GetAcceptLanguageHeader(r)

		// valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		slug, err := common.ReadSlugParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		catWithTrs, err := services.GetCategoryBySlugPublicService(app, slug, langCode)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		categoryPublicResponse := mappers.CategoryToCategoryPublicResponseMapper(catWithTrs.Category)

		detailResponse := types.NewDetailResponse(categoryPublicResponse, catWithTrs.Translations)
		err = common.WriteDetailJson(w, http.StatusOK, detailResponse, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

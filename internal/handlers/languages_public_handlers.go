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

// @Summary List languages
// @Description List languages with pagination and filters
// @Tags languages
// @Param Accept-Language header string false "Languages: en, ru, tk"
// @Param filters query requests.LanguagesAdminFilters true "Filters"
// @Produce json
// @Router /api/v1/languages [get]
// @Success 200 {object} types.PaginatedResponse[responses.LanguagePublicResponse]
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
func ListLanguagesPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		langCode := common.GetAcceptLanguageHeader(r)

		valTrans := r.Context().Value(constants.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		filters := requests.LanguagesAdminFilters{}

		qs := r.URL.Query()
		readLanguageAdminQueryParams(&filters, qs)
		err := app.Validator.Struct(&filters)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		langsWithTrs, metadata, err := services.ListLanguagesPublicService(app, &filters, langCode)
		langWithTransResponses := make([]*types.DetailResponse[responses.LanguagePublicResponse], len(langsWithTrs))
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		for _, langWithTrs := range langsWithTrs {
			langPublicResponse := mappers.LanguageToLanguagePublicResponseMapper(langWithTrs.Language)
			detailResponse := types.NewDetailResponse(langPublicResponse, langWithTrs.Translations)
			langWithTransResponses = append(langWithTransResponses, detailResponse)
		}

		// Write the response as JSON
		paginatedRes := types.PaginatedResponse[responses.LanguagePublicResponse]{
			Metadata: metadata,
			Results:  langWithTransResponses,
		}
		err = common.WritePaginatedJson(w, http.StatusOK, paginatedRes, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

// @Summary Get language by id
// @Description Get specific language details by id
// @Tags languages
// @Param Accept-Language header string false "Languages: en, ru, tk"
// @Param id path uuid true "UUID"
// @Produce json
// @Router /api/v1/languages/{id} [get]
// @Success 200 {object} types.DetailResponse[responses.LanguagePublicResponse]
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
func GetLanguagePublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(constants.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		language, err := services.GetLanguageService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		languageResponse := mappers.LanguageToLanguagePublicResponseMapper(language)

		detailResponse := types.NewDetailResponse(languageResponse, nil)
		err = common.WriteDetailJson(w, http.StatusOK, detailResponse, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

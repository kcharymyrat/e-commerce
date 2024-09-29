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
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func ListCategoriesPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		input := requests.ListCategoriesInput{}

		// Parse query string from the request
		qs := r.URL.Query()

		// Reading query parameters
		readCategoryQueryParameters(&input, qs)

		// Validate input
		err := app.Validator.Struct(input)
		if err != nil {
			errs := err.(validator.ValidationErrors)
			translatedErrs := make(map[string]string)
			for _, e := range errs {
				translatedErrs[e.Field()] = e.Translate(valTrans)
			}
			common.FailedValidationResponse(app.Logger, w, r, translatedErrs)
			return
		}

		// Retrieve categories from your data models
		categories, metadata, err := services.ListCategoriesService(app, input)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		// TODO: Map categories to category responses
		categoryPublicResponses := make([]*responses.CategoryPublicResponse, len(categories))
		for _, category := range categories {
			result := mappers.CategoryToCategoryPublicResponseMapper(category)
			categoryPublicResponses = append(categoryPublicResponses, result)
		}

		// Write the response as JSON
		err = common.WriteJson(w, http.StatusOK, common.Envelope{
			"metadata": metadata,
			"results":  categoryPublicResponses,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func GetCategoryPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		category, err := services.GetCategoryService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		// FIXME: Is it redundant?
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

		categoryPublicResponse := mappers.CategoryToCategoryPublicResponseMapper(category)

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": categoryPublicResponse}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

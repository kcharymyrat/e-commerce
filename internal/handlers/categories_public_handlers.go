package handlers

import (
	"errors"
	"net/http"

	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

func ListCategoriesPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		input := requests.ListCategoriesInput{}
		v := validator.New()

		// Parse query string from the request
		qs := r.URL.Query()

		// Reading query parameters
		readAndValidateQueryParameters(&input, qs, v)

		// Validate input using your filters
		filtersValidation(&input, v)

		// If validation fails, return a validation error response
		if !v.Valid() {
			common.FailedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		// Retrieve categories from your data models
		categories, metadata, err := services.ListCategoriesService(app, input)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
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
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func GetCategoryPublicHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, w, r)
			return
		}

		category, err := services.GetCategoryService(app, id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, w, r)
			default:
				common.ServerErrorResponse(app.Logger, w, r, err)
			}
			return
		}

		categoryPublicResponse := mappers.CategoryToCategoryPublicResponseMapper(category)

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": categoryPublicResponse}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

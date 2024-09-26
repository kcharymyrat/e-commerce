package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/mappers"
	"github.com/kcharymyrat/e-commerce/internal/services"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

func ListCategoriesManagerHandler(app *app.Application) http.HandlerFunc {
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
		categoryManagerResponses := make([]*responses.CategoryManagerResponse, len(categories))
		for _, category := range categories {
			result := mappers.CategoryToCategoryManagerResponseMapper(category)
			categoryManagerResponses = append(categoryManagerResponses, result)
		}

		// Write the response as JSON
		err = common.WriteJson(w, http.StatusOK, common.Envelope{
			"metadata": metadata,
			"results":  categoryManagerResponses,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func CreateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categoryInput requests.CreateCategoryInput

		err := common.ReadJSON(w, r, &categoryInput)
		if err != nil {
			common.BadRequestResponse(app.Logger, w, r, err)
			return
		}

		category := mappers.CreateCategoryInputToCategoryMapper(&categoryInput)

		v := validator.New()
		if data.ValidateCategory(v, category); !v.Valid() {
			common.FailedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		err = services.CreateCategoryService(app, category)
		if err != nil {
			HandleCategoryServiceErrors(w, r, app, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/api/v1/%d", category.ID))

		categoryResponse := mappers.CategoryToCategoryManagerResponseMapper(category)

		err = common.WriteJson(w, http.StatusCreated, common.Envelope{"category": categoryResponse}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func GetCategoryManagerHandler(app *app.Application) http.HandlerFunc {
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

		categoryManagerResponse := mappers.CategoryToCategoryManagerResponseMapper(category)

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": categoryManagerResponse}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func UpdateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
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
				return
			}
		}

		var input struct {
			ParentID    *uuid.UUID `json:"parent_id"`
			Name        string     `json:"name"`
			Slug        string     `json:"slug"`
			ImageUrl    string     `json:"image_url"`
			Description *string    `json:"description"`
			UpdatedByID uuid.UUID  `json:"created_by_id"`
		}

		err = common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, w, r, err)
			return
		}

		category.ParentID = input.ParentID
		category.Name = input.Name
		category.Slug = input.Slug
		category.ImageUrl = input.ImageUrl
		category.UpdatedByID = input.UpdatedByID
		category.Description = input.Description

		v := validator.New()

		if data.ValidateCategory(v, category); !v.Valid() {
			common.FailedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		err = app.Repositories.Categories.Update(category)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": category}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func PartialUpdateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, w, r)
			return
		}

		category, err := app.Repositories.Categories.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, w, r)
			default:
				common.ServerErrorResponse(app.Logger, w, r, err)
				return
			}
		}

		var input struct {
			Name        *string    `json:"name"`
			Slug        *string    `json:"slug"`
			ImageUrl    *string    `json:"image_url"`
			Description *string    `json:""`
			UpdatedByID *uuid.UUID `json:"created_by_id"`
		}

		err = common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, w, r, err)
			return
		}

		if input.Name != nil {
			category.Name = *input.Name
		}

		if input.Slug != nil {
			category.Slug = *input.Slug
		}

		if input.ImageUrl != nil {
			category.ImageUrl = *input.ImageUrl
		}

		if input.UpdatedByID != nil {
			message := "UpdatedById must be always provided if there is an update"
			common.ErrorResponse(app.Logger, w, r, http.StatusBadRequest, message)
			return
		}

		category.UpdatedByID = *input.UpdatedByID

		v := validator.New()

		if data.ValidateCategory(v, category); !v.Valid() {
			common.FailedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		err = app.Repositories.Categories.Update(category)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": category}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func DeleteCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, w, r)
			return
		}

		err = app.Repositories.Categories.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, w, r)
			default:
				common.ServerErrorResponse(app.Logger, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"message": "category successfully deleted"}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

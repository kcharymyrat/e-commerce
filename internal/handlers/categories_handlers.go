package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/common"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/filters"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

func ListCategoriesHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Names     []string    `json:"names"`
			Slugs     []string    `json:"slugs"`
			ParentIDs []uuid.UUID `json:"parent_ids"`
			filters.SearchFilters
			filters.CreatedUpdatedAtFilters
			filters.CreatedUpdatedByFilters
			filters.SortListFilters
			filters.PaginationFilters
		}

		v := validator.New()

		// Parse query string from the request
		qs := r.URL.Query()

		// Reading query parameters
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

		// Validate input using your filters
		filters.ValidateSearchFilters(v, input.SearchFilters)
		filters.ValidateCreatedUpdatedAtFilters(v, input.CreatedUpdatedAtFilters)
		filters.ValidateCreatedUpdatedByFilters(v, input.CreatedUpdatedByFilters)
		filters.ValidateSortFilters(v, input.SortListFilters)
		filters.ValidatePaginationFilters(v, input.PaginationFilters)

		// If validation fails, return a validation error response
		if !v.Valid() {
			common.FailedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		// Retrieve categories from your data models
		categories, metadata, err := app.Models.Categories.GetAll(
			input.Names,
			input.Slugs,
			input.ParentIDs,
			input.Search,
			input.CreatedAtFrom,
			input.CreatedAtUpTo,
			input.UpdatedAtFrom,
			input.UpdatedAtUpTo,
			input.CreatedByIDs,
			input.UpdatedByIDs,
			input.Sorts,
			input.SortSafeList,
			input.Page,
			input.PageSize,
		)

		// Handle any error that occurs during the database query
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
			return
		}

		// Write the response as JSON
		err = common.WriteJson(w, http.StatusOK, common.Envelope{
			"metadata": metadata,
			"results":  categories,
		}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func CreateCategoryHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categoryInput requests.CreateCategoryInput

		err := common.ReadJSON(w, r, &categoryInput)
		if err != nil {
			common.BadRequestResponse(app.Logger, w, r, err)
			return
		}

		category := &data.Category{
			ParentID:    categoryInput.ParentID,
			Name:        categoryInput.Name,
			Slug:        categoryInput.Slug,
			Description: categoryInput.Description,
			ImageUrl:    categoryInput.ImageUrl,
			CreatedByID: categoryInput.CreatedByID,
			UpdatedByID: categoryInput.UpdatedByID,
		}

		v := validator.New()

		if data.ValidateCategory(v, category); !v.Valid() {
			common.FailedValidationResponse(app.Logger, w, r, v.Errors)
			return
		}

		err = app.Models.Categories.Insert(category)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/api/v1/%d", category.ID))

		// TODO: map to category response

		err = common.WriteJson(w, http.StatusCreated, common.Envelope{"category": category}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func GetCategoryHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, w, r)
			return
		}

		category, err := app.Models.Categories.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, w, r)
			default:
				common.ServerErrorResponse(app.Logger, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": category}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, w, r, err)
		}
	}
}

func UpdateCategoryHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, w, r)
			return
		}

		category, err := app.Models.Categories.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
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

		err = app.Models.Categories.Update(category)
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

func PartialUpdateCategoryHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, w, r)
			return
		}

		category, err := app.Models.Categories.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
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

		err = app.Models.Categories.Update(category)
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

func DeleteCategoryHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, w, r)
			return
		}

		err = app.Models.Categories.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
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

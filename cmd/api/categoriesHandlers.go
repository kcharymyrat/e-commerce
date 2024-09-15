package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/filters"
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

func (app *application) listCategoriesHandler(w http.ResponseWriter, r *http.Request) {
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

	qs := r.URL.Query()

	input.Names = app.readQueryCSStrs(qs, "names")
	input.Slugs = app.readQueryCSStrs(qs, "slugs")
	input.ParentIDs = app.readQueryCSUUIDs(qs, "parent_ids", v)
	input.Search = app.readQueryStr(qs, "search")
	input.CreatedAtFrom = app.readQueryTime(qs, "created_at_from", v)
	input.CreatedAtUpTo = app.readQueryTime(qs, "created_at_up_to", v)
	input.UpdatedAtFrom = app.readQueryTime(qs, "updated_at_from", v)
	input.UpdatedAtUpTo = app.readQueryTime(qs, "updated_at_up_to", v)
	input.CreatedByIDs = app.readQueryCSUUIDs(qs, "created_by_ids", v)
	input.UpdatedByIDs = app.readQueryCSUUIDs(qs, "updated_by_ids", v)
	input.Sorts = app.readQueryCSStrs(qs, "sorts")
	input.SortSafeList = []string{
		"id", "name", "created_at", "updated_at", "-id", "-name", "-created_at", "-updated_at",
	}
	input.Page = app.readQueryInt(qs, "page", v)
	input.PageSize = app.readQueryInt(qs, "page_size", v)

	filters.ValidateSearchFilters(v, input.SearchFilters)
	filters.ValidateCreatedUpdatedAtFilters(v, input.CreatedUpdatedAtFilters)
	filters.ValidateCreatedUpdatedByFilters(v, input.CreatedUpdatedByFilters)
	filters.ValidateSortFilters(v, input.SortListFilters)
	filters.ValidatePaginationFilters(v, input.PaginationFilters)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

}

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var categoryInput struct {
		ParentId    uuid.UUID `json:"parent_id,omitempty"`
		Name        string    `json:"name"`
		Slug        string    `json:"slug"`
		Description string    `json:"description,omitempty"`
		ImageUrl    string    `json:"image_url"`
		CreatedByID uuid.UUID `json:"created_by_id"`
		UpdatedByID uuid.UUID `json:"updated_by_id"`
	}

	err := app.readJSON(w, r, &categoryInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category := &data.Category{
		ParentID:    categoryInput.ParentId,
		Name:        categoryInput.Name,
		Slug:        categoryInput.Slug,
		Description: categoryInput.Description,
		ImageUrl:    categoryInput.ImageUrl,
		CreatedByID: categoryInput.CreatedByID,
		UpdatedByID: categoryInput.UpdatedByID,
	}

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Insert(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/%d", category.ID))

	err = app.writeJson(w, http.StatusCreated, envelope{"category": category}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	var input struct {
		Name        string    `json:"name"`
		Slug        string    `json:"slug"`
		ImageUrl    string    `json:"image_url"`
		UpdatedByID uuid.UUID `json:"created_by_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category.Name = input.Name
	category.Slug = input.Slug
	category.ImageUrl = input.ImageUrl
	category.UpdatedByID = input.UpdatedByID

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Update(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) partialUpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	category, err := app.models.Categories.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	var input struct {
		Name        *string    `json:"name"`
		Slug        *string    `json:"slug"`
		ImageUrl    *string    `json:"image_url"`
		UpdatedByID *uuid.UUID `json:"created_by_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
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
		app.errorResponse(w, r, http.StatusBadRequest, message)
		return
	}

	category.UpdatedByID = *input.UpdatedByID

	v := validator.New()

	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Update(category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUUIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Categories.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelope{"message": "category successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

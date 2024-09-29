package handlers

import (
	"errors"
	"fmt"
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

func ListCategoriesManagerHandler(app *app.Application) http.HandlerFunc {
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
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func CreateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		var categoryInput requests.CreateCategoryInput
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
			HandleCategoryServiceErrors(app.Logger, localizer, w, r, err)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/api/v1/%d", category.ID))

		categoryResponse := mappers.CategoryToCategoryManagerResponseMapper(category)

		err = common.WriteJson(w, http.StatusCreated, common.Envelope{"category": categoryResponse}, headers)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func GetCategoryManagerHandler(app *app.Application) http.HandlerFunc {
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

		categoryManagerResponse := mappers.CategoryToCategoryManagerResponseMapper(category)

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": categoryManagerResponse}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func UpdateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
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
				return
			}
		}

		input := requests.UpdateCategoryInput{}

		err = common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
			return
		}

		category.ParentID = input.ParentID
		category.Name = input.Name
		category.Slug = input.Slug
		category.ImageUrl = input.ImageUrl
		category.UpdatedByID = input.UpdatedByID
		category.Description = input.Description

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

		err = app.Repositories.Categories.Update(category)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": category}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func PartialUpdateCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		category, err := app.Repositories.Categories.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
				return
			}
		}

		input := requests.PartialUpdateCategoryInput{}

		err = common.ReadJSON(w, r, &input)
		if err != nil {
			common.BadRequestResponse(app.Logger, localizer, w, r, err)
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

		category.UpdatedByID = input.UpdatedByID

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

		err = app.Repositories.Categories.Update(category)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"category": category}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

func DeleteCategoryManagerHandler(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// valTrans := r.Context().Value(common.ValTransKey).(ut.Translator)
		localizer := r.Context().Value(common.LocalizerKey).(*i18n.Localizer)

		id, err := common.ReadUUIDParam(r)
		if err != nil {
			common.NotFoundResponse(app.Logger, localizer, w, r)
			return
		}

		err = app.Repositories.Categories.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrRecordNotFound):
				common.NotFoundResponse(app.Logger, localizer, w, r)
			default:
				common.ServerErrorResponse(app.Logger, localizer, w, r, err)
			}
			return
		}

		err = common.WriteJson(w, http.StatusOK, common.Envelope{"message": "category successfully deleted"}, nil)
		if err != nil {
			common.ServerErrorResponse(app.Logger, localizer, w, r, err)
		}
	}
}

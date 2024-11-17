package services

import (
	"github.com/google/uuid"
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/types"
)

func CreateCategoryService(app *app.Application, category *data.Category) error {
	return app.Repositories.Categories.Create(category)
}

func GetCategoryByIDService(app *app.Application, id uuid.UUID) (*data.Category, error) {
	return app.Repositories.Categories.GetByID(id)
}

func GetCategoryBySlugService(app *app.Application, slug string) (*data.Category, error) {
	return app.Repositories.Categories.GetBySlug(slug)
}

func ListCategoriesService(
	app *app.Application,
	filters *requests.CategoriesAdminFilters,
) ([]*data.Category, types.Metadata, error) {
	return app.Repositories.Categories.List(filters)
}

func UpdateCategoryService(
	app *app.Application,
	input *requests.CategoryAdminUpdate,
	category *data.Category,
) error {
	category.ParentID = input.ParentID
	category.Name = input.Name
	category.Slug = input.Slug
	category.Description = input.Description
	category.ImageUrl = input.ImageUrl
	category.UpdatedByID = input.UpdatedByID

	return app.Repositories.Categories.Update(category)
}

func PartialUpdateCategoryService(
	app *app.Application,
	input *requests.CategoryAdminPartialUpdate,
	category *data.Category,
) error {
	if input.Name != nil {
		category.Name = *input.Name
	}

	if input.Slug != nil {
		category.Slug = *input.Slug
	}

	if input.Description != nil {
		category.Description = input.Description
	}

	if input.ImageUrl != nil {
		category.ImageUrl = *input.ImageUrl
	}

	category.UpdatedByID = input.UpdatedByID

	return app.Repositories.Categories.Update(category)
}

func DeleteCategoryServiceById(app *app.Application, id uuid.UUID) error {
	return app.Repositories.Categories.DeleteByID(id)
}

func DeleteCategoryServiceBySlug(app *app.Application, slug string) error {
	return app.Repositories.Categories.DeleteBySlug(slug)
}

package mappers

import (
	"github.com/kcharymyrat/e-commerce/api/requests"
	"github.com/kcharymyrat/e-commerce/api/responses"
	"github.com/kcharymyrat/e-commerce/internal/data"
)

func CreateCategoryInputToCategoryMapper(input *requests.CreateCategoryInput) *data.Category {
	return &data.Category{
		ParentID:    input.ParentID,
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		ImageUrl:    input.ImageUrl,
		CreatedByID: input.CreatedByID,
		UpdatedByID: input.UpdatedByID,
	}
}

func CategoryToCategoryPublicResponseMapper(category *data.Category) *responses.CategoryPublicResponse {
	// TODO:
	return nil
}

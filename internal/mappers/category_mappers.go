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
	return &responses.CategoryPublicResponse{
		ID:          category.ID,
		ParentID:    category.ParentID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageUrl:    category.ImageUrl,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func CategoryToCategoryManagerResponseMapper(category *data.Category) *responses.CategoryManagerResponse {
	return &responses.CategoryManagerResponse{
		ID:          category.ID,
		ParentID:    category.ParentID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ImageUrl:    category.ImageUrl,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
		CreatedByID: category.CreatedByID,
		UpdatedByID: category.UpdatedByID,
	}
}

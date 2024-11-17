package common

import (
	"math"

	"github.com/kcharymyrat/e-commerce/internal/types"
)

func CalculateMetadata(totalRecords, page, pageSize int) types.PaginationMetadata {
	if totalRecords <= 0 {
		return types.PaginationMetadata{}
	}

	return types.PaginationMetadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

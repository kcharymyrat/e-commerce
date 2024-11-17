package common

import (
	"math"

	"github.com/kcharymyrat/e-commerce/internal/types"
)

func CalculateMetadata(totalRecords, page, pageSize int) types.Metadata {
	if totalRecords <= 0 {
		return types.Metadata{}
	}

	return types.Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

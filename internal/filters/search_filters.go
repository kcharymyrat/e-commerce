package filters

import (
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

type SearchFilters struct {
	Search *string `json:"search"`
}

func ValidateSearchFilters(v *validator.Validator, f SearchFilters) {
	v.Check(len([]rune(*f.Search)) <= 50, "search", "must be maximum of 50 characters")
}

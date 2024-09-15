package filters

import (
	"github.com/kcharymyrat/e-commerce/internal/validator"
)

type PaginationFilters struct {
	Page     *int
	PageSize *int
}

func ValidatePaginationFilters(v *validator.Validator, f PaginationFilters) {
	if f.Page != nil {
		v.Check(*f.Page > 0, "page", "must be greater than zero")
		v.Check(*f.Page <= 10_000_000, "page", "must be a maximum of 10 millon")
	}
	if f.PageSize != nil {
		v.Check(*f.PageSize > 0, "page_size", "must be greater than zero")
		v.Check(*f.PageSize <= 100, "page_size", "must be a maximum of 100")
	}
}

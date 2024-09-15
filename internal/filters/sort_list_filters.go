package filters

import (
	"fmt"

	"github.com/kcharymyrat/e-commerce/internal/validator"
)

type SortListFilters struct {
	Sorts        []string `json:"sorts"`
	SortSafeList []string `json:"sort_safe_list"`
}

func ValidateSortFilters(v *validator.Validator, f SortListFilters) {
	for _, sort := range f.Sorts {
		v.Check(
			validator.In(sort, f.SortSafeList...),
			"sorts",
			fmt.Sprintf("invalid sort value: %s", sort),
		)
	}
}

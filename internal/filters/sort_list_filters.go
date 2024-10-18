package filters

import (
	"fmt"
	"strings"
)

type SortListFilter struct {
	Sorts        []string `json:"sorts" validate:"omitempty,dive,max=50"`
	SortSafeList []string `json:"sort_safe_list" validate:"omitempty,dive,max=50"`
}

func AddSortListFilterToSQL(f *SortListFilter, query *string) {
	if len(f.Sorts) > 0 {
		*query += " ORDER BY"
		for _, sort := range f.Sorts {
			direction := "ASC"
			sortField := strings.TrimSpace(strings.ToLower(sort))
			if strings.HasPrefix(sort, "-") {
				direction = "DESC"
				sortField = strings.TrimPrefix(sort, "-")
			}
			*query += fmt.Sprintf(" %s %s,", sortField, direction)
		}
		*query += " id ASC"
	}
}

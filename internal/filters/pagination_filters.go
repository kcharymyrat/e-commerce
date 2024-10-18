package filters

import "fmt"

type PaginationFilter struct {
	Page     *int `json:"page,omitempty" validate:"omitempty,gte=1,lte=10_000_000"`
	PageSize *int `json:"page_size,omitempty" validate:"omitempty,gte=1,lte=100"`
}

func AddPaginationFilterToSQL(
	f *PaginationFilter, query *string, argCounter *int, args []interface{},
) {
	fallbackPageSize := 20 // FIXME: make a constant number
	if f.PageSize != nil {
		*query += fmt.Sprintf(" LIMIT $%d", *argCounter)
		args = append(args, *f.PageSize)
		*argCounter++
		fallbackPageSize = *f.PageSize
	} else {
		f.PageSize = &fallbackPageSize
		*query += fmt.Sprintf(" LIMIT %d", fallbackPageSize)
	}

	defaultPage := 1 // FIXME: make a constant number
	if f.Page != nil {
		offset := fallbackPageSize * (*f.Page - 1)
		*query += fmt.Sprintf(" OFFSET $%d", *argCounter)
		args = append(args, offset)
		*argCounter++
	} else {
		f.Page = &defaultPage
		*query += fmt.Sprintf(" OFFSET %d", *f.Page)
	}
}

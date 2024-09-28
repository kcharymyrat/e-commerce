package filters

type PaginationFilters struct {
	Page     *int `json:"page,omitempty" validate:"omitempty,gte=1,lte=10_000_000"`
	PageSize *int `json:"page_size,omitempty" validate:"omitempty,gte=1,lte=100"`
}

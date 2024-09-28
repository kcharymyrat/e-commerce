package filters

type SortListFilters struct {
	Sorts        []string `json:"sorts" validate:"omitempty,dive,max=50"`
	SortSafeList []string `json:"sort_safe_list" validate:"omitempty,dive,max=50"`
}

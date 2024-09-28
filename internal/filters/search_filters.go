package filters

type SearchFilters struct {
	Search *string `json:"search" validate:"omitempty,max=50"`
}

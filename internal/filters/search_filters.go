package filters

type SearchFilter struct {
	Search *string `json:"search" validate:"omitempty,max=50"`
}

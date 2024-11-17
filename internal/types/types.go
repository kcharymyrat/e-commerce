package types

type ContextKey string

type UserClaimsKey struct{}

type Envelope map[string]interface{}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

type PaginatedResults struct {
	Metadata interface{} `json:"metadata"`
	Results  interface{} `json:"results"`
}

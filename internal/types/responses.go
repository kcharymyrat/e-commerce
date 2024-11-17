package types

type ErrorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type DetailResponse[T any] struct {
	Data *T `json:"data"`
}

type PaginatedResponse[T any] struct {
	Metadata PaginationMetadata `json:"metadata"`
	Results  []*T               `json:"results"`
}

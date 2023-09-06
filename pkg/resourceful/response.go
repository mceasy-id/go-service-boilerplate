package resourceful

type Response[IDType, Model any] struct {
	Metadata *Metadata           `json:"metadata"`
	Data     Data[IDType, Model] `json:"data"`
}

type Metadata struct {
	Count      int `json:"count"`
	Page       int `json:"page"`
	TotalCount int `json:"total_count"`
	TotalPage  int `json:"total_page"`
}

type Data[IDType, Model any] struct {
	PaginatedResult []Model  `json:"paginated_result"`
	Ids             []IDType `json:"ids"`
}

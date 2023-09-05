package entities

type ResponseData struct {
	Data   any                 `json:"data"`
	Errors map[string][]string `json:"errors,omitempty"`
}

type GetProductOption struct {
	PessimisticLocking bool
}

package dtos

import "github.com/invopop/validation"

type IndexRequest struct {
	Limit   int      `json:"limit"`
	Page    int      `json:"page"`
	Search  string   `json:"search"`
	Filters []string `json:"filters"`
	Sorts   []string `json:"sort"`
}

func (i IndexRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Limit, validation.Required, validation.Min(1)),
		validation.Field(&i.Page, validation.Required, validation.Min(1)),
	)
}

package dtos

import (
	"strings"

	"github.com/invopop/validation"
)

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func (c *CreateProductRequest) Mod() *CreateProductRequest {
	c.Name = strings.Join(strings.Fields(c.Name), " ")
	c.Description = strings.Join(strings.Fields(c.Description), " ")

	return c
}

func (c CreateProductRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(0, 25)),
	)
}

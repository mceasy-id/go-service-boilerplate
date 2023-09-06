package dtos

import (
	"strings"

	"github.com/invopop/validation"
)

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func (u *UpdateProductRequest) Mod() *UpdateProductRequest {
	u.Name = strings.Join(strings.Fields(u.Name), " ")
	u.Description = strings.Join(strings.Fields(u.Description), " ")
	return u
}

func (u UpdateProductRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(0, 25)),
	)
}

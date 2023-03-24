package models

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Price       *int64    `json:"price,omitempty"`
}

type StoreProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       *int64 `json:"price,omitempty"`
}

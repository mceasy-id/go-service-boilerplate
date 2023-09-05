package dtos

import "mceasy/service-demo/internal/product/entities"

type ShowProductResponse struct {
	UUID        string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
}

func NewProductResponse(entity entities.Product) ShowProductResponse {
	return ShowProductResponse{
		UUID:        entity.UUID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
	}
}

package product

import (
	"context"
	"mceasy/service-demo/internal/models"
)

type Repository interface {
	GetProducts(ctx context.Context) ([]*models.Product, error)
	StoreProduct(ctx context.Context, product *models.Product) (*models.Product, error)
}

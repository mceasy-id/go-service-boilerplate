package product

import (
	"context"
	"mceasy/service-demo/internal/models"
)

type UseCase interface {
	Index(ctx context.Context) ([]*models.Product, error)
	Store(ctx context.Context, product *models.Product) (*models.Product, error)
}

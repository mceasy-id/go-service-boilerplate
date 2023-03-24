package usecase

import (
	"context"
	"mceasy/service-demo/internal/models"
	"mceasy/service-demo/internal/product"
)

type productUC struct {
	productRepo product.Repository
}

func NewProductUseCase(productRepo product.Repository) product.UseCase {
	return &productUC{
		productRepo: productRepo,
	}
}

func (u *productUC) Index(ctx context.Context) ([]*models.Product, error) {
	products, err := u.productRepo.GetProducts(ctx)
	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		products = make([]*models.Product, 0)
	}

	return products, nil
}

func (u *productUC) Store(ctx context.Context, product *models.Product) (*models.Product, error) {
	result, err := u.productRepo.StoreProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	return result, nil
}

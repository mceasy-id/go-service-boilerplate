package repository

import (
	"context"
	"database/sql"
	"mceasy/service-demo/internal/models"
	"mceasy/service-demo/internal/product"

	"github.com/google/uuid"
)

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) product.Repository {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) GetProducts(ctx context.Context) ([]*models.Product, error) {
	query := getProducts

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var products []*models.Product
	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price)
		if err != nil {
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (r *productRepo) StoreProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	newUUID := uuid.New()

	query := storeProduct
	_, err := r.db.ExecContext(ctx,
		query,
		newUUID,
		product.Name,
		product.Description,
		product.Price,
	)

	if err != nil {
		return nil, err
	}

	var createdProduct models.Product
	createdProduct.ID = newUUID

	return &createdProduct, nil
}

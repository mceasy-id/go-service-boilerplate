package product

import (
	"context"
	"database/sql"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	"mceasy/service-demo/pkg/resourceful"

	"github.com/google/uuid"
)

type Repository interface {
	Atomic(ctx context.Context, opt *sql.TxOptions, cb func(tx Repository) error) error

	GetProductByUUID(ctx context.Context, productUUID string, companyId int64, options ...entities.GetProductOption) (entities.Product, error)
	StoreNewProduct(ctx context.Context, product entities.Product) (string, error)
	DeleteProductByUUID(ctx context.Context, productUUID string, companyId int64) error
	UpdateProductByUUID(ctx context.Context, product entities.UpdateProduct) error

	FindProductResourceful(ctx context.Context, resource *resourceful.Resource[uuid.UUID, dtos.ProductList]) (*resourceful.Resource[uuid.UUID, dtos.ProductList], error)
	IsProductKeyExists(ctx context.Context, payload entities.StoreProduct, companyId int64) (bool, error)
}

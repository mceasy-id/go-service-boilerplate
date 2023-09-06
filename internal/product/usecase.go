package product

import (
	"context"
	"mceasy/service-demo/internal/identity/identityentities"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	"mceasy/service-demo/pkg/resourceful"

	"github.com/google/uuid"
)

type UseCase interface {
	Index(ctx context.Context, companyId uint64, resource *resourceful.Resource[uuid.UUID, dtos.ProductList]) (*resourceful.Resource[uuid.UUID, dtos.ProductList], error)
	Store(ctx context.Context, requestCredential identityentities.Credential, payload entities.StoreProduct) (string, error)
	Show(ctx context.Context, productUUID string, companyId int64) (entities.Product, error)
	Update(ctx context.Context, requestCredential identityentities.Credential, productUUID string, payload entities.UpdateProduct) error
	Delete(ctx context.Context, productUUID string, companyId int64) error
}

package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mceasy/service-demo/internal/identity/identityentities"
	"mceasy/service-demo/internal/product"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	"mceasy/service-demo/pkg/apperror"
	"mceasy/service-demo/pkg/observability/instrumentation"
	"mceasy/service-demo/pkg/resourceful"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UseCaseParameter struct {
	ProductRepo product.Repository
}

func NewProductUseCase(param UseCaseParameter) product.UseCase {
	return &productUC{
		repo: param.ProductRepo,
	}
}

type productUC struct {
	repo product.Repository
}

func (u *productUC) Update(ctx context.Context, requestCredential identityentities.Credential, productUUID string, payload entities.UpdateProduct) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"UpdateUseCase",
	)
	defer span.End()

	return u.repo.Atomic(ctx, &sql.TxOptions{}, func(tx product.Repository) error {
		productEntity, err := tx.GetProductByUUID(ctx, productUUID, int64(requestCredential.CompanyId), entities.GetProductOption{
			PessimisticLocking: true,
		})
		if err != nil {
			return err
		}

		exists, err := tx.IsProductKeyExists(ctx, entities.StoreProduct{
			Name:        payload.Name,
			Description: payload.Description,
			Price:       productEntity.Price,
		}, int64(requestCredential.CompanyId))
		if err != nil {
			return err
		}

		if !strings.EqualFold(productEntity.Name, payload.Name) && exists {
			return apperror.BadRequestMap(map[string][]string{
				"product": {"already exists"},
			})
		}

		err = tx.UpdateProductByUUID(ctx, entities.UpdateProduct{
			CompanyId:   int64(requestCredential.CompanyId),
			UUID:        productUUID,
			Name:        payload.Name,
			Description: payload.Description,
			Price:       payload.Price,
			UpdatedOn:   time.Now(),
			UpdatedBy:   requestCredential.UserName,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (u *productUC) Delete(ctx context.Context, productUUID string, companyId int64) error {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"DeleteUseCase",
	)
	defer span.End()

	err := u.repo.DeleteProductByUUID(ctx, productUUID, companyId)
	if err != nil {
		return err
	}

	return nil
}

func (u *productUC) Index(ctx context.Context, companyId uint64, resource *resourceful.Resource[uuid.UUID, dtos.ProductList]) (*resourceful.Resource[uuid.UUID, dtos.ProductList], error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"IndexUseCase",
	)
	defer span.End()

	resource.Parameter.LocalFilters = append(
		resource.Parameter.LocalFilters,
		fmt.Sprintf("company_id eq %d", companyId),
	)

	err := resource.SetParam(*resource.Parameter)
	if err != nil {
		return nil, err
	}

	productResource, err := u.repo.FindProductResourceful(ctx, resource)
	if err != nil {
		if errors.Is(err, resourceful.ErrPagination) {
			return resource, nil
		}
		return nil, err
	}

	return productResource, nil
}

func (u *productUC) Store(ctx context.Context, requestCredential identityentities.Credential, payload entities.StoreProduct) (string, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"StoreUseCase",
	)
	defer span.End()

	exists, err := u.repo.IsProductKeyExists(ctx, payload, int64(requestCredential.CompanyId))
	if err != nil {
		return "", err
	}

	if exists {
		return "", apperror.BadRequestMap(map[string][]string{
			"product": {"already exists"},
		})
	}

	productUUID, err := u.repo.StoreNewProduct(ctx, entities.NewProduct(requestCredential, payload))
	if err != nil {
		return "", err
	}

	return productUUID, nil
}

func (uc *productUC) Show(ctx context.Context, productUUID string, companyId int64) (entities.Product, error) {
	ctx, span := instrumentation.NewTraceSpan(
		ctx,
		"ShowUseCase",
	)
	defer span.End()

	product, err := uc.repo.GetProductByUUID(ctx, productUUID, companyId)
	if err != nil {
		return entities.Product{}, err
	}

	return product, nil
}

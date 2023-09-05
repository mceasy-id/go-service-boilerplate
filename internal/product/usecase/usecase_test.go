package usecase_test

import (
	"context"
	"fmt"
	"mceasy/service-demo/internal/identity/identityentities"
	v1 "mceasy/service-demo/internal/product/delivery/http/external/v1"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	mocks "mceasy/service-demo/internal/product/mock"
	"mceasy/service-demo/internal/product/usecase"
	"mceasy/service-demo/pkg/observability/instrumentation"
	"mceasy/service-demo/pkg/resourceful"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func init() {
	v1.NewProductInstance()
}

func TestProductUseCase_Show(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedCtx, _ := instrumentation.NewTraceSpan(
		ctx,
		"ShowUseCase",
	)

	expectedProduct := entities.Product{
		CompanyId: 392,
		UUID:      uuid.NewString(),
	}

	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockProductRepo.
		EXPECT().
		GetProductByUUID(
			expectedCtx,
			expectedProduct.UUID,
			expectedProduct.CompanyId,
		).Return(expectedProduct, nil)

	productUC := usecase.NewProductUseCase(
		usecase.UseCaseParameter{
			ProductRepo: mockProductRepo,
		},
	)

	_, err := productUC.Show(ctx, expectedProduct.UUID, expectedProduct.CompanyId)
	require.NoError(t, err)
}

func TestProductUseCase_Index(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedCtx, _ := instrumentation.NewTraceSpan(
		ctx,
		"IndexUseCase",
	)

	expectedResourceful := resourceful.NewResource[uuid.UUID, dtos.ProductList](v1.ProductDefinition)
	err := expectedResourceful.SetParam(resourceful.Parameter{
		Limit: 10,
		Page:  1,
		LocalFilters: []string{
			"company_id eq 392",
		},
	})
	require.NoError(t, err)

	returnedResourceful := resourceful.NewResource[uuid.UUID, dtos.ProductList](v1.ProductDefinition)
	err = returnedResourceful.SetParam(resourceful.Parameter{
		Limit: 10,
		Page:  1,
		LocalFilters: []string{
			"company_id eq 392",
		},
	})
	require.NoError(t, err)

	returnedResourceful.SetResult(resourceful.Result[uuid.UUID, dtos.ProductList]{
		Ids: []uuid.UUID{uuid.New()},
		PaginatedResult: []dtos.ProductList{
			{
				Name:        "Kacang",
				Description: "Ini Kacang",
			},
		},
	})

	mockProductRepo := mocks.NewMockRepository(ctrl)

	mockProductRepo.
		EXPECT().
		FindProductResourceful(expectedCtx, expectedResourceful).
		Return(returnedResourceful, nil)

	productUC := usecase.NewProductUseCase(
		usecase.UseCaseParameter{
			ProductRepo: mockProductRepo,
		},
	)

	resourceProduct := resourceful.NewResource[uuid.UUID, dtos.ProductList](v1.ProductDefinition)
	err = resourceProduct.SetParam(resourceful.Parameter{
		Limit: 10,
		Page:  1,
	})
	require.NoError(t, err)

	resourcefulProduct, err := productUC.Index(ctx, 392, resourceProduct)
	require.NoError(t, err)
	assert.NotEqual(t, resourcefulProduct, resourceProduct)
}

func TestProductUseCase_Store(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedCtx, _ := instrumentation.NewTraceSpan(
		ctx,
		"StoreUseCase",
	)

	expectedCred := identityentities.Credential{
		UserName:  "admin",
		UserId:    1,
		CompanyId: 392,
	}

	expectedUUID := uuid.NewString()

	storeProduct := entities.StoreProduct{
		Name:        "Mie Indomie",
		Description: "Mi enak!",
		Price:       2500,
	}

	mockProductRepo := mocks.NewMockRepository(ctrl)

	mockProductRepo.
		EXPECT().
		IsProductKeyExists(
			expectedCtx,
			storeProduct,
			int64(392),
		).Return(false, nil)

	mockProductRepo.
		EXPECT().
		StoreNewProduct(
			expectedCtx,
			createProductMatcher(entities.Product{
				CompanyId:   392,
				Name:        storeProduct.Name,
				Description: storeProduct.Description,
				Price:       storeProduct.Price,
				CreatedBy:   expectedCred.UserName,
			}),
		).Return(expectedUUID, nil)

	productUC := usecase.NewProductUseCase(usecase.UseCaseParameter{
		ProductRepo: mockProductRepo,
	})

	_, err := productUC.Store(ctx, expectedCred, storeProduct)
	require.NoError(t, err)
}

func TestProductUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	expectedCtx, _ := instrumentation.NewTraceSpan(
		ctx,
		"DeleteUseCase",
	)

	expectedProductUUID := uuid.NewString()

	mockProductRepo := mocks.NewMockRepository(ctrl)

	mockProductRepo.
		EXPECT().
		DeleteProductByUUID(
			expectedCtx,
			expectedProductUUID,
			int64(392),
		).Return(nil)

	productUC := usecase.NewProductUseCase(usecase.UseCaseParameter{
		ProductRepo: mockProductRepo,
	})

	err := productUC.Delete(ctx, expectedProductUUID, 392)
	require.NoError(t, err)
}

func createProductMatcher(product entities.Product) gomock.Matcher {
	return eqProductMatcher{
		product: product,
	}
}

type eqProductMatcher struct {
	product entities.Product
}

func (e eqProductMatcher) Matches(x interface{}) bool {
	arg, ok := x.(entities.Product)
	if !ok {
		return false
	}

	return arg.CompanyId == e.product.CompanyId &&
		arg.Name == e.product.Name &&
		arg.Description == e.product.Description &&
		arg.Price == e.product.Price &&
		arg.CreatedBy == e.product.CreatedBy
}

func (e eqProductMatcher) String() string {
	return fmt.Sprintf("%v", e.product.Name)
}

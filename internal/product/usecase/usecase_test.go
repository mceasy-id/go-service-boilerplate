package usecase

import (
	"context"
	"mceasy/service-demo/internal/models"
	"mceasy/service-demo/internal/product/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Index(t *testing.T) {
	// Acquire
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	// Expect Mock
	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockProducts := []*models.Product{{}}
	mockProductRepo.EXPECT().GetProducts(ctx).Return(mockProducts, nil)

	// Action
	productUC := NewProductUseCase(mockProductRepo)
	result, err := productUC.Index(ctx)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

func Test_Store(t *testing.T) {
	// Acquire
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	// Expect Input
	inputProductPrice := int64(12000)
	inputProduct := &models.Product{
		Name:        "Kecap Manis",
		Description: "Kecap rasa manis",
		Price:       &inputProductPrice,
	}

	// Expect Mock
	mockProductRepo := mocks.NewMockRepository(ctrl)
	mockProduct := &models.Product{ID: uuid.New()}
	mockProductRepo.EXPECT().StoreProduct(ctx, inputProduct).Return(mockProduct, nil)

	// Action
	productUC := NewProductUseCase(mockProductRepo)
	result, err := productUC.Store(ctx, inputProduct)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, result)
}

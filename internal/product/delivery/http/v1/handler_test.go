package v1

import (
	"mceasy/service-demo/internal/models"
	"mceasy/service-demo/internal/product/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Index(t *testing.T) {
	// Acquire
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fiberApp := fiber.New()

	// Expect Mock
	mockProductUC := mocks.NewMockUseCase(ctrl)
	products := []*models.Product{{}}
	mockProductUC.EXPECT().Index(gomock.Any()).Return(products, nil)

	// Action
	productHandler := NewProductHandler(mockProductUC)
	MapProductHandlers(fiberApp, productHandler)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response, err := fiberApp.Test(request)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func Test_Store(t *testing.T) {
	// Acquire
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	fiberApp := fiber.New()

	// Expect Input
	inputJson := `{
		"name": "Indomie",
		"description": "Indomie Rebus",
		"price": 10000
	}`

	var inputProduct models.Product
	inputProduct.Name = "Indomie"
	inputProduct.Description = "Indomie Rebus"
	inputProductPrice := int64(10000)
	inputProduct.Price = &inputProductPrice

	// Expect Mock
	mockProductUC := mocks.NewMockUseCase(ctrl)
	product := models.Product{ID: uuid.New()}
	mockProductUC.EXPECT().Store(gomock.Any(), &inputProduct).Return(&product, nil)

	// Action
	productHandler := NewProductHandler(mockProductUC)
	MapProductHandlers(fiberApp, productHandler)

	request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(inputJson))
	request.Header.Set("Content-Type", "application/json")
	response, err := fiberApp.Test(request)

	// Assert
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

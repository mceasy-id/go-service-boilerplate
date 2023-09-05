package v1_test

import (
	"encoding/json"
	"fmt"
	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/identity/identityentities"
	v1 "mceasy/service-demo/internal/product/delivery/http/external/v1"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	mocks "mceasy/service-demo/internal/product/mock"
	"mceasy/service-demo/pkg/apperror"
	"mceasy/service-demo/pkg/resourceful"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newFiberApp(companyId uint64) *fiber.App {
	var fiberConfig fiber.Config
	fiberConfig.ErrorHandler = apperror.HttpHandleError
	fiberApp := fiber.New(fiberConfig)

	fiberApp.Use(func(c *fiber.Ctx) error {
		c.Locals("authCredential", identityentities.Credential{
			CompanyId: companyId,
			UserName:  "cavalry",
			UserId:    1,
		})

		return c.Next()
	})

	return fiberApp

}

func init() {
	v1.NewProductInstance()
}

func TestProductHandler_Show(t *testing.T) {
	t.Run("contract_test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedProductUUID := uuid.NewString()
		mockProductUC := mocks.NewMockUseCase(ctrl)

		returnProductEntity := entities.Product{
			UUID:        expectedProductUUID,
			Name:        "Mie Sedap !",
			Description: "Sedap sekali",
			Price:       2000,
		}

		mockProductUC.
			EXPECT().
			Show(
				gomock.Any(),
				expectedProductUUID,
				int64(392),
			).Return(returnProductEntity, nil)

		productHandler := v1.NewProductHandler(config.Config{}, mockProductUC)
		app := newFiberApp(392)
		v1.MapProduct(app, productHandler)

		req := httptest.NewRequest(
			http.MethodGet,
			fmt.Sprintf("/product/%s", expectedProductUUID),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var contract struct {
			Data struct {
				Id          string `json:"id"`
				Name        string `json:"name"`
				Description string `json:"description"`
				Price       uint64 `json:"price"`
			}
		}

		jsonDecoder := json.NewDecoder(resp.Body)
		jsonDecoder.DisallowUnknownFields()

		err = jsonDecoder.Decode(&contract)
		require.NoError(t, err)

		require.Equal(t, contract.Data.Id, returnProductEntity.UUID)
		require.Equal(t, contract.Data.Name, returnProductEntity.Name)
		require.Equal(t, contract.Data.Description, returnProductEntity.Description)
	})
}

func TestProductHandler_Index(t *testing.T) {
	t.Run("contract_test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		queryParameter := "?limit=10&page=1&sorts=name+desc&search=Mie"

		resourceParam := resourceful.Parameter{
			Limit:  10,
			Page:   1,
			Search: "Mie",
			Sorts:  []string{"name desc"},
		}

		expectedResource := resourceful.NewResource[uuid.UUID, dtos.ProductList](v1.ProductDefinition)

		expectedResource.Parameter = &resourceParam

		length_product := 5
		ids := make([]uuid.UUID, 0, length_product)
		products := make([]dtos.ProductList, 0, length_product)

		for i := 0; i < length_product; i++ {
			productUUID := uuid.New()
			ids = append(ids, productUUID)

			products = append(products, dtos.ProductList{
				UUID:        productUUID.String(),
				Name:        faker.WORD,
				Description: faker.WORD,
			})
		}

		returnedResource := resourceful.NewResource[uuid.UUID, dtos.ProductList](v1.ProductDefinition)
		err := returnedResource.SetParam(resourceParam)
		require.NoError(t, err)

		returnedResource.SetResult(resourceful.Result[uuid.UUID, dtos.ProductList]{
			Ids:             ids,
			PaginatedResult: products,
		})

		// Call mock usecase
		productUCMock := mocks.NewMockUseCase(ctrl)
		productUCMock.EXPECT().Index(gomock.Any(), uint64(392), expectedResource).Return(returnedResource, nil)

		// Action
		productHandler := v1.NewProductHandler(config.Config{}, productUCMock)
		app := newFiberApp(392)
		v1.MapProduct(app, productHandler)

		request := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/product%s", queryParameter), nil)
		response, err := app.Test(request)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, response.StatusCode)

		var contract struct {
			Metadata struct {
				Count      int `json:"count"`
				Page       int `json:"page"`
				TotalPage  int `json:"total_page"`
				TotalCount int `json:"total_count"`
			} `json:"metadata"`
			Data struct {
				PaginatedResult []struct {
					Id          *string `json:"id"`
					Name        *string `json:"name"`
					Description *string `json:"description"`
				} `json:"paginated_result"`
				Ids []string `json:"ids"`
			} `json:"data"`
		}

		jsonDecoder := json.NewDecoder(response.Body)
		jsonDecoder.DisallowUnknownFields()
		err = jsonDecoder.Decode(&contract)
		require.NoError(t, err)

		assert.Equal(t, length_product, contract.Metadata.Count)
		assert.Equal(t, 1, contract.Metadata.Page)
		assert.Equal(t, 1, contract.Metadata.TotalPage)
		assert.Equal(t, 5, contract.Metadata.Count)
		assert.NotEmpty(t, contract.Data.Ids)
		assert.NotEmpty(t, contract.Data.Ids[0])
		assert.NotEmpty(t, contract.Data.PaginatedResult[0].Id)
		assert.NotEmpty(t, contract.Data.PaginatedResult[0].Name)
		assert.NotEmpty(t, contract.Data.PaginatedResult[0].Description)
	})
}

func TestProductHandler_Store(t *testing.T) {
	t.Run("valid_store", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		requestInStruct := dtos.CreateProductRequest{
			Name:        "Sunlight",
			Description: "Buat Cuci Piring",
			Price:       2000,
		}
		requestInJSON, err := json.Marshal(requestInStruct)
		require.NoError(t, err)

		productUUID := uuid.NewString()

		mockProductUC := mocks.NewMockUseCase(ctrl)
		mockProductUC.
			EXPECT().
			Store(
				gomock.Any(),
				identityentities.Credential{
					UserName:  "cavalry",
					UserId:    1,
					CompanyId: 392,
				},
				entities.StoreProduct{
					Name:        "Sunlight",
					Description: "Buat Cuci Piring",
					Price:       2000,
				},
			).Return(productUUID, nil)

		productHandler := v1.NewProductHandler(config.Config{}, mockProductUC)
		app := newFiberApp(392)
		v1.MapProduct(app, productHandler)

		req := httptest.NewRequest(
			http.MethodPost,
			"/product",
			strings.NewReader(string(requestInJSON)),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestProductHandler_Delete(t *testing.T) {
	t.Run("Ok_when_valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		productUUID := uuid.NewString()

		productUCMock := mocks.NewMockUseCase(ctrl)

		productUCMock.EXPECT().Delete(
			gomock.Any(),
			productUUID,
			int64(392),
		).Return(nil)

		reasonHandler := v1.NewProductHandler(config.Config{}, productUCMock)

		app := newFiberApp(392)
		v1.MapProduct(app, reasonHandler)

		req := httptest.NewRequest(
			http.MethodDelete,
			fmt.Sprintf("/product/%s", productUUID),
			nil,
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestProductHandler_Update(t *testing.T) {
	t.Run("ok_when_valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		requestInStruct := dtos.UpdateProductRequest{
			Name:        "Rinso",
			Description: "Ini rinso",
			Price:       2000,
		}
		requestInJSON, err := json.Marshal(requestInStruct)
		require.NoError(t, err)

		productUUID := uuid.NewString()

		mockProductUC := mocks.NewMockUseCase(ctrl)
		mockProductUC.
			EXPECT().
			Update(
				gomock.Any(),
				identityentities.Credential{
					UserName:  "cavalry",
					UserId:    1,
					CompanyId: 392,
				},
				productUUID,
				updateProductMatcher(entities.UpdateProduct{
					CompanyId:   392,
					UUID:        productUUID,
					Name:        requestInStruct.Name,
					Description: requestInStruct.Description,
					Price:       requestInStruct.Price,
					UpdatedOn:   time.Now(),
					UpdatedBy:   "cavalry",
				}),
			).Return(nil)

		productHandler := v1.NewProductHandler(config.Config{}, mockProductUC)
		app := newFiberApp(392)
		v1.MapProduct(app, productHandler)

		req := httptest.NewRequest(
			http.MethodPatch,
			fmt.Sprintf("/product/%s", productUUID),
			strings.NewReader(string(requestInJSON)),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

	})
}

func updateProductMatcher(product entities.UpdateProduct) gomock.Matcher {
	return eqUpdateProductMatcher{
		product: product,
	}
}

type eqUpdateProductMatcher struct {
	product entities.UpdateProduct
}

func (e eqUpdateProductMatcher) Matches(x interface{}) bool {
	arg, ok := x.(entities.UpdateProduct)
	if !ok {
		return false
	}

	return arg.CompanyId == e.product.CompanyId &&
		arg.Name == e.product.Name &&
		arg.Description == e.product.Description &&
		arg.Price == e.product.Price &&
		arg.UpdatedBy == e.product.UpdatedBy
}

func (e eqUpdateProductMatcher) String() string {
	return fmt.Sprintf("%v", e.product)
}

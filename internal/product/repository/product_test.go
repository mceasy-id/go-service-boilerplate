package repository_test

import (
	"context"
	"database/sql"
	"log"
	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/product"
	v1 "mceasy/service-demo/internal/product/delivery/http/external/v1"
	"mceasy/service-demo/internal/product/dtos"
	"mceasy/service-demo/internal/product/entities"
	"mceasy/service-demo/internal/product/repository"
	"mceasy/service-demo/pkg/database"
	"mceasy/service-demo/pkg/resourceful"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var cfg config.Config

func init() {
	loadCfg, err := config.LoadConfigPath("../../../config/config-local")
	if err != nil {
		log.Fatal(err)
	}

	cfg = loadCfg

	v1.NewProductInstance()
}

func TestProductRepo_Atomic(t *testing.T) {
	db, err := database.GetPostgreConnection(cfg)
	require.NoError(t, err)
	productRepo := repository.NewProductPGRepo(db)

	t.Run("ok_and_commit", func(t *testing.T) {
		err := productRepo.Atomic(context.TODO(), &sql.TxOptions{}, func(tx product.Repository) error {
			expectedUUID := uuid.NewString()

			_, err := tx.StoreNewProduct(context.Background(), entities.Product{
				CompanyId:   392,
				UUID:        expectedUUID,
				Name:        "Kecap Bango",
				Description: "Ini kecap",
				Price:       1000,
				CreatedOn:   time.Now(),
				CreatedBy:   "admin",
				UpdatedOn:   time.Now(),
				UpdatedBy:   "admin",
			})
			if err != nil {
				return err
			}

			_, err = tx.GetProductByUUID(context.Background(), expectedUUID, 392)
			if err != nil {
				return err
			}

			return nil
		})

		require.NoError(t, err)
	})

	t.Run("when_error_should_rollback", func(t *testing.T) {
		err := productRepo.Atomic(context.TODO(), &sql.TxOptions{}, func(tx product.Repository) error {
			expectedUUID := uuid.NewString()

			_, err := tx.StoreNewProduct(context.Background(), entities.Product{
				CompanyId:   392,
				UUID:        expectedUUID,
				Name:        "Garam Dapur",
				Description: "Ini garam",
				Price:       2000,
				CreatedOn:   time.Now(),
				CreatedBy:   "admin",
				UpdatedOn:   time.Now(),
				UpdatedBy:   "admin",
			})
			if err != nil {
				return err
			}

			_, err = tx.GetProductByUUID(context.Background(), uuid.NewString(), 392)
			if err != nil {
				return err
			}

			return nil
		})
		require.Error(t, err)
	})

}

func TestProductRepo_GetProductByUUID(t *testing.T) {
	db, err := database.GetPostgreConnection(cfg)
	require.NoError(t, err)
	productRepo := repository.NewProductPGRepo(db)

	t.Run("integration_not_empty", func(t *testing.T) {
		expectedProduct := entities.Product{
			CompanyId:   392,
			UUID:        uuid.NewString(),
			Name:        "Daia",
			Description: "ini deterjen",
			Price:       1000,
			CreatedOn:   time.Now(),
			CreatedBy:   "admin",
			UpdatedOn:   time.Now(),
			UpdatedBy:   "admin",
		}

		_, err := productRepo.StoreNewProduct(context.Background(), expectedProduct)
		require.NoError(t, err)

		product, err := productRepo.GetProductByUUID(context.Background(), expectedProduct.UUID, expectedProduct.CompanyId)
		require.NoError(t, err)

		require.Equal(t, expectedProduct.CompanyId, product.CompanyId)
		require.Equal(t, expectedProduct.Name, product.Name)
		require.Equal(t, expectedProduct.Description, product.Description)
		require.Equal(t, expectedProduct.Price, product.Price)
		require.Equal(t, expectedProduct.CreatedBy, product.CreatedBy)
		require.Equal(t, expectedProduct.UpdatedBy, product.UpdatedBy)
	})

	t.Run("integration_not_empty_with_option", func(t *testing.T) {
		expectedProduct := entities.Product{
			CompanyId:   392,
			UUID:        uuid.NewString(),
			Name:        "Daia",
			Description: "ini deterjen",
			Price:       1000,
			CreatedOn:   time.Now(),
			CreatedBy:   "admin",
			UpdatedOn:   time.Now(),
			UpdatedBy:   "admin",
		}

		_, err := productRepo.StoreNewProduct(context.Background(), expectedProduct)
		require.NoError(t, err)

		product, err := productRepo.GetProductByUUID(context.Background(),
			expectedProduct.UUID,
			expectedProduct.CompanyId,
			entities.GetProductOption{
				PessimisticLocking: true,
			},
		)
		require.NoError(t, err)

		require.Equal(t, expectedProduct.CompanyId, product.CompanyId)
		require.Equal(t, expectedProduct.Name, product.Name)
		require.Equal(t, expectedProduct.Description, product.Description)
		require.Equal(t, expectedProduct.Price, product.Price)
		require.Equal(t, expectedProduct.CreatedBy, product.CreatedBy)
		require.Equal(t, expectedProduct.UpdatedBy, product.UpdatedBy)
	})

}

func TestProductRepo_FindProductResourceful(t *testing.T) {
	db, err := database.GetPostgreConnection(cfg)
	require.NoError(t, err)
	productRepo := repository.NewProductPGRepo(db)

	t.Run("integration_not_empty", func(t *testing.T) {
		expectedProduct := entities.Product{
			CompanyId:   392,
			UUID:        uuid.NewString(),
			Name:        "Daia",
			Description: "ini deterjen",
			Price:       1000,
			CreatedOn:   time.Now(),
			CreatedBy:   "admin",
			UpdatedOn:   time.Now(),
			UpdatedBy:   "admin",
		}

		_, err := productRepo.StoreNewProduct(context.Background(), expectedProduct)
		require.NoError(t, err)

		instance := resourceful.NewResource[uuid.UUID, dtos.ProductList](v1.ProductDefinition)
		err = instance.SetParam(resourceful.Parameter{Limit: 10, Page: 1, LocalFilters: []string{"company_id eq 392"}})
		require.NoError(t, err)

		resourceProduct, err := productRepo.FindProductResourceful(context.Background(), instance)
		require.NoError(t, err)
		require.NotEmpty(t, resourceProduct)

		resourceProductResponse := resourceProduct.Response()
		require.NotEmpty(t, resourceProductResponse)
		require.NotEmpty(t, resourceProductResponse.Data)
		require.NotEmpty(t, resourceProductResponse.Metadata)
		require.NotEmpty(t, resourceProductResponse.Data.Ids)
		require.NotEmpty(t, resourceProductResponse.Data.PaginatedResult)
		require.NotEmpty(t, resourceProductResponse.Data.PaginatedResult[0].UUID)
		require.NotEmpty(t, resourceProductResponse.Data.PaginatedResult[0].Name)
		require.NotEmpty(t, resourceProductResponse.Data.PaginatedResult[0].Description)
	})

}

func TestProductRepo_DeleteProductByUUID(t *testing.T) {
	db, err := database.GetPostgreConnection(cfg)
	require.NoError(t, err)
	productRepo := repository.NewProductPGRepo(db)

	t.Run("integration_ok", func(t *testing.T) {
		expectedProduct := entities.Product{
			CompanyId:   392,
			UUID:        uuid.NewString(),
			Name:        "Swallow",
			Description: "Ini Sandal",
			Price:       5000,
			CreatedOn:   time.Now(),
			CreatedBy:   "admin",
			UpdatedOn:   time.Now(),
			UpdatedBy:   "admin",
		}

		_, err := productRepo.StoreNewProduct(context.Background(), expectedProduct)
		require.NoError(t, err)

		err = productRepo.DeleteProductByUUID(context.Background(), expectedProduct.UUID, expectedProduct.CompanyId)
		require.NoError(t, err)

		_, err = productRepo.GetProductByUUID(context.Background(), expectedProduct.UUID, expectedProduct.CompanyId)
		require.Error(t, err)
	})
}

func TestProductRepo_UpdateProductByUUID(t *testing.T) {
	db, err := database.GetPostgreConnection(cfg)
	require.NoError(t, err)
	productRepo := repository.NewProductPGRepo(db)

	t.Run("integration_ok", func(t *testing.T) {
		expectedProduct := entities.Product{
			CompanyId:   392,
			UUID:        uuid.NewString(),
			Name:        "Swallow",
			Description: "Ini Sandal",
			Price:       5000,
			CreatedOn:   time.Now(),
			CreatedBy:   "admin",
			UpdatedOn:   time.Now(),
			UpdatedBy:   "admin",
		}

		_, err := productRepo.StoreNewProduct(context.Background(), expectedProduct)
		require.NoError(t, err)

		err = productRepo.UpdateProductByUUID(
			context.Background(),
			entities.UpdateProduct{
				CompanyId:   expectedProduct.CompanyId,
				UUID:        expectedProduct.UUID,
				Name:        "BATA",
				Description: "Sandal Bata",
				Price:       10000,
				UpdatedOn:   time.Now(),
				UpdatedBy:   "who !?",
			},
		)
		require.NoError(t, err)

		product, err := productRepo.GetProductByUUID(
			context.Background(),
			expectedProduct.UUID,
			expectedProduct.CompanyId,
		)
		require.NoError(t, err)

		require.Equal(t, expectedProduct.CompanyId, product.CompanyId)
		require.Equal(t, expectedProduct.UUID, product.UUID)
		require.Equal(t, "BATA", product.Name)
		require.Equal(t, "Sandal Bata", product.Description)
		require.Equal(t, int64(10000), product.Price)
		require.Equal(t, expectedProduct.CreatedBy, product.CreatedBy)
		require.Equal(t, "who !?", product.UpdatedBy)
	})
}

package repository

import (
	"context"
	"database/sql"
	"log"
	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/models"
	"mceasy/service-demo/pkg/database"
	"testing"

	"github.com/stretchr/testify/require"
)

var dbConn *sql.DB

func init() {
	cfg, err := config.LoadConfigPath("../../../config/config")
	if err != nil {
		log.Fatal(err)
	}

	dbConn, _ = database.GetDatabaseConnection(&cfg)
}

func Test_GetProducts(t *testing.T) {
	if dbConn == nil {
		t.Skip()
	}

	productRepo := NewProductRepository(dbConn)

	// Action
	_, err := productRepo.GetProducts(context.Background())

	// Assert
	require.NoError(t, err)
}

func Test_StoreProduct(t *testing.T) {
	if dbConn == nil {
		t.Skip()
	}

	productRepo := NewProductRepository(dbConn)

	var product models.Product
	product.Name = "Kecap Asin"
	product.Description = "Kecap rasa asin"
	productPrice := int64(12000)
	product.Price = &productPrice

	// Action
	createdProduct, err := productRepo.StoreProduct(context.Background(), &product)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, createdProduct)

	products, _ := productRepo.GetProducts(context.Background())
	require.NotEmpty(t, products)
}

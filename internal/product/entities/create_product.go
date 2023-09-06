package entities

import (
	"mceasy/service-demo/internal/identity/identityentities"
	"time"

	"github.com/google/uuid"
)

type StoreProduct struct {
	Name        string
	Description string
	Price       int64
}

func NewProduct(cred identityentities.Credential, storeProduct StoreProduct) Product {
	return Product{
		CompanyId:   int64(cred.CompanyId),
		UUID:        uuid.NewString(),
		Name:        storeProduct.Name,
		Description: storeProduct.Description,
		Price:       storeProduct.Price,
		CreatedOn:   time.Now(),
		CreatedBy:   cred.UserName,
		UpdatedOn:   time.Now(),
		UpdatedBy:   cred.UserName,
	}
}

package v1

import (
	"log"
	"mceasy/service-demo/internal/product/tabledefinition"
	"mceasy/service-demo/pkg/resourceful"
)

var (
	ProductDefinition *resourceful.Definition
)

func NewProductInstance() {
	productDefinition, err := resourceful.NewDefinition(tabledefinition.Product)
	if err != nil {
		log.Println(err)
	}

	ProductDefinition = productDefinition
}

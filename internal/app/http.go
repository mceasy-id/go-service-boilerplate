package app

import (
	productHttp "mceasy/service-demo/internal/product/delivery/http/v1"
	productRepository "mceasy/service-demo/internal/product/repository"
	productUseCase "mceasy/service-demo/internal/product/usecase"
)

func (a *app) MapHttpHandlers() {
	productRepo := productRepository.NewProductRepository(a.DB)
	productUC := productUseCase.NewProductUseCase(productRepo)
	productHandler := productHttp.NewProductHandler(productUC)

	v1 := a.Fiber.Group("api/v1")
	productV1 := v1.Group("product")

	productHttp.MapProductHandlers(productV1, productHandler)
}

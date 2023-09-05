package server

import (
	productHttpV1 "mceasy/service-demo/internal/product/delivery/http/external/v1"
	productRepository "mceasy/service-demo/internal/product/repository"
	productUseCase "mceasy/service-demo/internal/product/usecase"

	"mceasy/service-demo/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) MapHandlers() error {
	check := s.Fiber.Group("/check")
	check.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	// FE to service
	webV1 := s.Fiber.Group("/external/api/web/v1")
	// Service to service
	internalV1 := s.Fiber.Group("/internal/api/v1")

	// CORS Middleware
	webV1.Use(middleware.CORSMiddleware(s.Config.App))

	// JWT Middleware
	webV1.Use(middleware.GuardMiddleware(s.Config))

	// Internal Middleware
	internalV1.Use(middleware.InternalMiddleware(s.Config))

	// Resource Initialization
	productHttpV1.NewProductInstance()

	//* App repository - Internal
	productRepo := productRepository.NewProductPGRepo(s.DB)

	//* App Use Case
	productUC := productUseCase.NewProductUseCase(
		productUseCase.UseCaseParameter{
			ProductRepo: productRepo,
		},
	)

	//* App Handler
	// productHandler := productHttpV1.NewProductHandler(s.Config, productUC)
	productHandler := productHttpV1.NewProductHandler(s.Config, productUC)
	productHttpV1.MapProduct(webV1, productHandler)

	return nil
}

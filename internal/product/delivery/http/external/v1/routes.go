package v1

import (
	"mceasy/service-demo/internal/product"

	"github.com/gofiber/fiber/v2"
)

func MapProduct(routes fiber.Router, h product.Handlers) {
	product := routes.Group("/product")
	product.Get("/", h.Index)
	product.Post("/", h.Store)
	product.Get("/:productUUID", h.Show)
	product.Delete("/:productId", h.Delete)
	product.Patch("/:productId", h.Update)
}

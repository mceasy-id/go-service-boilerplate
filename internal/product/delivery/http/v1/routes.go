package v1

import (
	"mceasy/service-demo/internal/product"

	"github.com/gofiber/fiber/v2"
)

func MapProductHandlers(routes fiber.Router, h product.Handlers) {
	routes.Get("", h.Index)
	routes.Post("", h.Store)
}

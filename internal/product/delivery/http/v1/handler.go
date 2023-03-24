package v1

import (
	"mceasy/service-demo/internal/models"
	"mceasy/service-demo/internal/product"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type productHandler struct {
	productUC product.UseCase
}

func NewProductHandler(productUC product.UseCase) product.Handlers {
	return &productHandler{
		productUC: productUC,
	}
}

func (h *productHandler) Index(c *fiber.Ctx) error {
	products, err := h.productUC.Index(c.UserContext())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.JSON(products)
}

func (h *productHandler) Store(c *fiber.Ctx) error {
	var request models.StoreProductRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	var product models.Product
	product.Name = request.Name
	product.Description = request.Description
	product.Price = request.Price

	result, err := h.productUC.Store(c.UserContext(), &product)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	return c.JSON(result)
}

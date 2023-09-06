package product

import "github.com/gofiber/fiber/v2"

type Handlers interface {
	Index(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

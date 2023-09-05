package middleware

import (
	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/identity/identityentities"
	"mceasy/service-demo/pkg/apperror"

	"github.com/gofiber/fiber/v2"
)

func InternalMiddleware(cfg config.Config) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		token, err := getTokenFromHeader(ctx)
		if err != nil {
			return err
		}

		if token != cfg.App.Key {
			return apperror.Unauthorized()
		}

		ctx.Locals("authCredential", identityentities.Credential{})

		return ctx.Next()
	}
}

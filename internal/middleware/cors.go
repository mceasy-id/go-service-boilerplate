package middleware

import (
	"strings"

	"mceasy/service-demo/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware(appConfig config.AppConfig) func(*fiber.Ctx) error {
	var allowedOrigins []string
	switch appConfig.Env {
	case "local", "dev":
		allowedOrigins = append(allowedOrigins, "http://localhost:8080", "https://*.mceasy.com")
	case "staging", "production":
		allowedOrigins = append(allowedOrigins, "https://*.mceasy.com")
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(allowedOrigins, ","),
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: true,
		AllowHeaders:     "Authorization,Content-Type,Traceparent",
	})
}

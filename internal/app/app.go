package app

import (
	"database/sql"
	"fmt"
	"log"
	"mceasy/service-demo/config"
	"mceasy/service-demo/pkg/database"

	"github.com/gofiber/fiber/v2"
)

type app struct {
	Config *config.Config
	Fiber  *fiber.App
	DB     *sql.DB
}

func NewApp(config *config.Config) *app {
	// Get DB Connection
	db, err := database.GetDatabaseConnection(config)
	if err != nil {
		log.Fatalf("Error NewApp() get database connection: %v", err)
	}

	// New Fiber
	fiberApp := fiber.New()

	return &app{
		Config: config,
		Fiber:  fiberApp,
		DB:     db,
	}
}

func (a *app) Run() error {
	// Map Http Handlers
	a.MapHttpHandlers()

	return a.Fiber.Listen(fmt.Sprintf(":%s", a.Config.App.Port))
}

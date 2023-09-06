package main

import (
	"context"
	"log"
	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/server"
	"mceasy/service-demo/pkg/database"
	"mceasy/service-demo/pkg/logger"
	"mceasy/service-demo/pkg/observability"
	"os"
)

func main() {
	log.Println("Starting App")

	// Get Environtment Config
	env := os.Getenv("config")
	if env == "" {
		env = "local"
	}

	// Load YAML Config
	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	//Init Logger
	appLogger, err := logger.InitLogger(cfg)
	if err != nil {
		log.Fatalf("Error on initializing logger: %s", err)
	}
	defer appLogger.Sync()

	// Init DB Connection
	db, err := database.GetPostgreConnection(cfg)
	if err != nil {
		appLogger.Fatal("Error on getting database postgre connection: %s", err)
	}

	if cfg.Observability.Enable {
		tracerProvider, err := observability.InitTracerProvider(&cfg)
		if err != nil {
			appLogger.Fatalf("Error on initializing tracer provider: %s", err)
		}
		defer tracerProvider.Shutdown(context.Background())

		meterProvider, err := observability.InitMeterProvider(&cfg)
		if err != nil {
			appLogger.Fatalf("Error on initializing meter provider: %s", err)
		}
		defer meterProvider.Shutdown(context.Background())
	}

	// Create & Run Server
	server := server.NewServer(cfg, appLogger, db)
	if err = server.Run(); err != nil {
		log.Fatal(err)
	}

	appLogger.Info("App stopped")

}

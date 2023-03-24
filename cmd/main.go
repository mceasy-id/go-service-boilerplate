package main

import (
	"log"
	"mceasy/service-demo/config"
	"mceasy/service-demo/internal/app"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app := app.NewApp(&cfg)
	if err := app.Run(); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}

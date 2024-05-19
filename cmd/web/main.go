package main

import (
	"crowdfunding-service/internal/config"
	"fmt"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	objectStore := config.NewMinio(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validate:    validate,
		Config:      viperConfig,
		ObjectStore: objectStore,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"context"
	"log"

	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/config"
	"gitlab.com/s0j0hn/go-rest-boilerplate-echo/internal/app"

	// Import docs for swagger generation
	_ "gitlab.com/s0j0hn/go-rest-boilerplate-echo/docs"
)

// @title Go REST API Boilerplate
// @version 2.0
// @description A modern REST API boilerplate built with Echo framework
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create application
	application := app.NewApp(cfg)

	// Initialize application
	if err := application.Initialize(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Start application
	ctx := context.Background()
	if err := application.Start(ctx); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}
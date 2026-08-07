package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/hoshina-dev/custapi/docs"
	"github.com/hoshina-dev/custapi/internal/config"
	"github.com/hoshina-dev/custapi/internal/database"
	"github.com/hoshina-dev/custapi/internal/handlers"
	"github.com/hoshina-dev/custapi/internal/repositories"
	"github.com/hoshina-dev/custapi/internal/routes"
	"github.com/hoshina-dev/custapi/internal/services"
	"github.com/hoshina-dev/custapi/internal/telemetry"
)

// @title				Customer API
// @version			1.0
// @description		A simple REST API for managing users and organizations
// @BasePath			/api/v1
//
// @tag.name			organizations
// @tag.description	Operations related to organizations
//
// @tag.name			users
// @tag.description	Operations related to users
func main() {
	// Load configuration
	cfg := config.Load()

	// Wire up tracing/metrics before anything else touches the network or
	// the database, so both are covered from the very first request.
	ctx := context.Background()
	shutdownTelemetry, err := telemetry.Setup(ctx, cfg.Telemetry)
	if err != nil {
		log.Fatalf("Failed to set up telemetry: %v", err)
	}

	// Initialize database
	db := database.ConnectDB(cfg.DataSourceName)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	orgRepo := repositories.NewOrganizationRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, orgRepo)
	orgService := services.NewOrganizationService(orgRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	orgHandler := handlers.NewOrgHandler(orgService)

	// Setup routes
	routes.SetupRoutes(app, userHandler, orgHandler, cfg.CorsOrigins)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		log.Printf("Starting server on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Failed to shutdown gracefully: %v", err)
	}

	// Flush any buffered spans/metrics before exiting so the final
	// requests of this process aren't silently dropped.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdownTelemetry(shutdownCtx); err != nil {
		log.Printf("Failed to shut down telemetry cleanly: %v", err)
	}

	log.Println("Server stopped")
}

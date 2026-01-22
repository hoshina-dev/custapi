package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	_ "github.com/hoshina-dev/custapi/docs"
	"github.com/hoshina-dev/custapi/internal/config"
	"github.com/hoshina-dev/custapi/internal/database"
	"github.com/hoshina-dev/custapi/internal/handlers"
	"github.com/hoshina-dev/custapi/internal/repositories"
	"github.com/hoshina-dev/custapi/internal/routes"
	"github.com/hoshina-dev/custapi/internal/services"
)

// @title			Customer API
// @version		1.0
// @description	A simple REST API for managing users and organizations
// @host			localhost:8080
// @BasePath		/api/v1
func main() {
	// Load configuration
	cfg := config.Load()

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
	routes.SetupRoutes(app, userHandler, orgHandler)

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

	log.Println("Server stopped")
}

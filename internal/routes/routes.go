package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/hoshina-dev/custapi/internal/handlers"
	"github.com/hoshina-dev/custapi/internal/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, orgHandler *handlers.OrgHandler) {
	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(middleware.Logger())
	app.Use(middleware.ErrorHandler())

	app.Get("/swagger/*", swagger.HandlerDefault)

	// API v1
	v1 := app.Group("/api/v1")
	{
		// Users routes
		v1.Get("/users", userHandler.GetUsers)
		v1.Get("/users/:id", userHandler.GetUser)
		v1.Get("/users/organization/:org_id", userHandler.GetUsersByOrganization)

		// Organizations routes
		v1.Get("/organizations", orgHandler.GetOrganizations)
		v1.Get("/organizations/:id", orgHandler.GetOrganization)
	}
}

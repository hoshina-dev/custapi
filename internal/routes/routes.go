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
		users := v1.Group("/users")
		users.Get("/", userHandler.GetUsers)
		users.Get("/search", userHandler.SearchUsers)
		users.Get("/organization/:org_id", userHandler.GetUsersByOrganization)
		users.Get("/:id", userHandler.GetUser)

		// Organizations routes
		org := v1.Group("/organizations")
		org.Get("/", orgHandler.GetOrganizations)
		org.Get("/search", orgHandler.SearchOrganizations)
		org.Get("/coordinates", orgHandler.GetAllCoords)
		org.Get("/:id", orgHandler.GetOrganization)
		org.Post("/", orgHandler.CreateOrganization)
		org.Post("/batch", orgHandler.GetByIDs)
		org.Patch("/:id", orgHandler.UpdateOrganization)
		org.Delete("/:id", orgHandler.DeleteOrganization)
	}
}

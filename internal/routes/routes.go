package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/hoshina-dev/custapi/internal/handlers"
	"github.com/hoshina-dev/custapi/internal/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, orgHandler *handlers.OrgHandler, healthHandler *handlers.HealthHandler, corsOrigins string) {
	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins,
	}))
	app.Use(middleware.Logger())
	app.Use(middleware.ErrorHandler())

	app.Get("/health", healthHandler.Health)
	app.Get("/ready", healthHandler.Ready)

	app.Get("/swagger/*", swagger.HandlerDefault)
	fmt.Println("📖 Swagger docs available at http://localhost:8080/swagger")

	// Scalar API Reference UI
	app.Get("/scalar", handlers.ScalarHandler)
	fmt.Println("📖 Scalar docs available at http://localhost:8080/scalar")

	// API v1
	v1 := app.Group("/api/v1")
	{
		// Users routes
		user := v1.Group("/users")
		user.Get("/", userHandler.GetUsers)
		user.Get("/search", userHandler.SearchUsers)
		user.Get("/id/:id", userHandler.GetUser)
		user.Get("/id/:id/organizations", userHandler.GetUserOrganizations)
		user.Get("/email/:email", userHandler.GetUserByEmail)
		user.Post("/", userHandler.CreateUser)
		user.Patch("/id/:id", userHandler.UpdateUser)
		user.Delete("/id/:id", userHandler.DeleteUser)

		// Auth routes
		auth := v1.Group("/auth")
		auth.Post("/verify", userHandler.VerifyCredentials)

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
		// Organization membership routes
		org.Get("/:id/members", orgHandler.GetMembers)
		org.Post("/:id/members", orgHandler.AddMember)
		org.Delete("/:id/members/:user_id", orgHandler.RemoveMember)
		org.Patch("/:id/members/:user_id", orgHandler.SetRole)
	}
}

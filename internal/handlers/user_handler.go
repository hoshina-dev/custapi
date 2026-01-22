package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hoshina-dev/custapi/internal/models"
	"github.com/hoshina-dev/custapi/internal/services"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Get a list of all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.UserResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userService.ListUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := make([]models.UserResponse, len(users))
	for i, u := range users {
		response[i] = models.UserResponse{
			ID:             u.ID,
			Email:          u.Email,
			Name:           u.Name,
			OrganizationID: u.OrganizationID,
		}
	}

	return c.JSON(response)
}

// GetUser godoc
//
//	@Summary		Get a user by ID
//	@Description	Get a single user by their ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	models.UserResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user, err := h.userService.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
	}

	response := models.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.Name,
		OrganizationID: user.OrganizationID,
	}

	return c.JSON(response)
}

// GetUsersByOrganization godoc
//
//	@Summary		Get users by organization
//	@Description	Get all users in a specific organization
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			org_id	path		string	true	"Organization ID"
//	@Success		200		{array}		models.UserResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users/organization/{org_id} [get]
func (h *UserHandler) GetUsersByOrganization(c *fiber.Ctx) error {
	orgID := c.Params("org_id")

	users, err := h.userService.ListUsersByOrganization(c.Context(), orgID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := make([]models.UserResponse, len(users))
	for i, u := range users {
		response[i] = models.UserResponse{
			ID:             u.ID,
			Email:          u.Email,
			Name:           u.Name,
			OrganizationID: u.OrganizationID,
		}
	}

	return c.JSON(response)
}

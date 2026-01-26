package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"github.com/hoshina-dev/custapi/internal/services"
)

// UserHandler handles user HTTP requests
type UserHandler struct {
	userService services.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// CreateUser godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user with email, name, organization, and optional details
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.CreateUserRequest	true	"User to create"
//	@Success		201		{object}	models.UserResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		422		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	req := new(models.CreateUserRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid json payload"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Error: err.Error()})
	}

	user, err := h.userService.CreateUser(c.Context(), req)
	if err != nil {
		if err.Error() == "organization not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user.ToResponse())
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
		response[i] = u.ToResponse()
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
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	user, err := h.userService.GetUser(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
	}

	return c.JSON(user.ToResponse())
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
	orgID, err := uuid.Parse(c.Params("org_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid organization id"})
	}

	users, err := h.userService.ListUsersByOrganization(c.Context(), orgID)
	if err != nil {
		if err.Error() == "organization not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := make([]models.UserResponse, len(users))
	for i, u := range users {
		response[i] = u.ToResponse()
	}

	return c.JSON(response)
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update an existing user by ID (partial updates supported)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string						true	"User ID (UUID)"
//	@Param			user	body		models.UpdateUserRequest	true	"Fields to update"
//	@Success		200		{object}	models.UserResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		422		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users/{id} [patch]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	req := new(models.UpdateUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid JSON payload"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Error: err.Error()})
	}

	user, err := h.userService.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
	}

	return c.JSON(user.ToResponse())
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Soft delete a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID (UUID)"
//	@Success		204
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	if err := h.userService.Delete(c.Context(), id); err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

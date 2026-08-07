package handlers

import (
	"net/mail"
	"net/url"

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

	user, err := h.userService.CreateUser(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "failed to create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(user.ToResponse())
}

// VerifyCredentials godoc
//
//	@Summary		Verify user credentials
//	@Description	Verify an email + password against the stored hash. Returns the user (without password) on success, 401 on failure. Used by the BFF for login so the password hash never leaves this service.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		models.VerifyCredentialsRequest	true	"Credentials to verify"
//	@Success		200			{object}	models.UserResponse
//	@Failure		400			{object}	models.ErrorResponse
//	@Failure		401			{object}	models.ErrorResponse
//	@Failure		422			{object}	models.ErrorResponse
//	@Failure		500			{object}	models.ErrorResponse
//	@Router			/auth/verify [post]
func (h *UserHandler) VerifyCredentials(c *fiber.Ctx) error {
	req := new(models.VerifyCredentialsRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid json payload"})
	}
	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Error: err.Error()})
	}
	user, err := h.userService.VerifyCredentials(c.UserContext(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "failed to verify credentials"})
	}
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{Error: "invalid email or password"})
	}
	return c.JSON(user.ToResponse())
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Returns a list of all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.UserResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userService.ListUsers(c.UserContext())
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
//	@Description	Get a single user by their UUID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID (UUID)"
//	@Success		200	{object}	models.UserDetailResponse
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/id/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	user, err := h.userService.GetUser(c.UserContext(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
	}

	return c.JSON(user.ToDetailResponse())
}

// GetUserByEmail godoc
//
//	@Summary		Get a user by email
//	@Description	Get a single user by their email address
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			email	path		string	true	"User email address"
//	@Success		200		{object}	models.UserDetailResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users/email/{email} [get]
func (h *UserHandler) GetUserByEmail(c *fiber.Ctx) error {
	email, err := url.PathUnescape(c.Params("email"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid email parameter"})
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid email format"})
	}

	user, err := h.userService.GetUserByEmail(c.UserContext(), email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "user not found"})
	}

	return c.JSON(user.ToDetailResponse())
}

// GetUserOrganizations godoc
//
//	@Summary		Get a user's organizations
//	@Description	Get all organizations a user belongs to with their role
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID (UUID)"
//	@Success		200	{array}		models.UserMembershipResponse
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/id/{id}/organizations [get]
func (h *UserHandler) GetUserOrganizations(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	memberships, err := h.userService.GetUserOrganizations(c.UserContext(), id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := make([]models.UserMembershipResponse, len(memberships))
	for i, m := range memberships {
		response[i] = m.ToMembershipResponse()
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
//	@Router			/users/id/{id} [patch]
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

	user, err := h.userService.Update(c.UserContext(), id, req)
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
//	@Router			/users/id/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid user id"})
	}

	if err := h.userService.Delete(c.UserContext(), id); err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// SearchUsers godoc
//
//	@Summary		Search users
//	@Description	Search users by name or email using ILIKE query
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			q		query		string	true	"Search query"
//	@Param			limit	query		int		false	"Maximum number of results to return (default: 100)"
//	@Success		200		{array}		models.UserResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users/search [get]
func (h *UserHandler) SearchUsers(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "query parameter 'q' is required"})
	}

	// Parse limit parameter with default of 100
	limit := c.QueryInt("limit", 100)
	if limit < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "limit must be non-negative"})
	}

	users, err := h.userService.SearchUsers(c.UserContext(), query, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := make([]models.UserResponse, len(users))
	for i, u := range users {
		response[i] = u.ToResponse()
	}

	return c.JSON(response)
}

package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"github.com/hoshina-dev/custapi/internal/services"
)

// OrgHandler handles organization HTTP requests
type OrgHandler struct {
	orgService services.OrganizationService
	validate   *validator.Validate
}

// NewOrgHandler creates a new organization handler
func NewOrgHandler(orgService services.OrganizationService) *OrgHandler {
	return &OrgHandler{
		orgService: orgService,
		validate:   validator.New(),
	}
}

// CreateOrganization godoc
// @Summary Create a new organization
// @Description Create a new organization with name, location, and optional details
// @Tags organizations
// @Accept json
// @Produce json
// @Param organization body models.CreateOrganizationRequest true "Organization to create"
// @Success 201 {object} models.OrganizationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations [post]
func (h *OrgHandler) CreateOrganization(c *fiber.Ctx) error {
	req := new(models.CreateOrganizationRequest)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid json payload"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Error: err.Error()})
	}

	org, err := h.orgService.CreateOrganization(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: "failed to create organization"})
	}

	return c.Status(fiber.StatusCreated).JSON(org.ToResponse())
}

// GetOrganizations godoc
// @Summary Get all organizations
// @Description Get a list of all organizations
// @Tags organizations
// @Accept json
// @Produce json
// @Success 200 {array} models.OrganizationResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations [get]
func (h *OrgHandler) GetOrganizations(c *fiber.Ctx) error {
	orgs, err := h.orgService.ListOrganizations(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	response := make([]models.OrganizationResponse, len(orgs))
	for i, o := range orgs {
		response[i] = o.ToResponse()
	}

	return c.JSON(response)
}

// GetOrganization godoc
// @Summary Get an organization by ID
// @Description Get a single organization by their ID
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID"
// @Success 200 {object} models.OrganizationResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/{id} [get]
func (h *OrgHandler) GetOrganization(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid organization id"})
	}

	org, err := h.orgService.GetOrganization(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if org == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "organization not found"})
	}

	return c.JSON(org.ToResponse())
}

// GetAllCoords godoc
// @Summary Get all organization coordinates
// @Description Get ID and coordinates of all organizations
// @Tags organizations
// @Accept json
// @Produce json
// @Success 200 {array} models.OrganizationCoord
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/coordinates [get]
func (h *OrgHandler) GetAllCoords(c *fiber.Ctx) error {
	orgs, err := h.orgService.GetAllCoords(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	coords := make([]models.OrganizationCoord, len(orgs))
	for i, org := range orgs {
		coords[i] = models.OrganizationCoord{ID: org.ID, Latitude: *org.Latitude, Longitude: *org.Longitude}
	}

	return c.JSON(coords)
}

// UpdateOrganization godoc
// @Summary Update an organization
// @Description Update an existing organization by ID (partial updates supported)
// @Tags organizations
// @Accept json
// @Produce json
// @Param id path string true "Organization ID (UUID)"
// @Param organization body models.UpdateOrganizationRequest true "Fields to update"
// @Success 200 {object} models.OrganizationResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 422 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /organizations/{id} [patch]
func (h *OrgHandler) UpdateOrganization(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid organization id"})
	}

	req := new(models.UpdateOrganizationRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid JSON payload"})
	}

	if err := h.validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.ErrorResponse{Error: err.Error()})
	}

	org, err := h.orgService.UpdateOrganization(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if org == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "organization not found"})
	}

	return c.JSON(org.ToResponse())
}

func (h *OrgHandler) DeleteOrganization(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid organization id"})
	}

	if err := h.orgService.DeleteOrganization(c.Context(), id); err != nil {
		if err.Error() == "organization not found" {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

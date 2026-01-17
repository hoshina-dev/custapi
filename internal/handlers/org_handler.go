package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hoshina-dev/custapi/internal/models"
	"github.com/hoshina-dev/custapi/internal/services"
)

// OrgHandler handles organization HTTP requests
type OrgHandler struct {
	orgService services.OrganizationService
}

// NewOrgHandler creates a new organization handler
func NewOrgHandler(orgService services.OrganizationService) *OrgHandler {
	return &OrgHandler{
		orgService: orgService,
	}
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
		response[i] = models.OrganizationResponse{
			ID:   o.ID.String(),
			Name: o.Name,
		}
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
	id := c.Params("id")

	org, err := h.orgService.GetOrganization(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{Error: err.Error()})
	}

	if org == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{Error: "organization not found"})
	}

	response := models.OrganizationResponse{
		ID:   org.ID.String(),
		Name: org.Name,
	}

	return c.JSON(response)
}

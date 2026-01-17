package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"github.com/hoshina-dev/custapi/internal/repositories"
)

// OrganizationService defines organization business logic operations
type OrganizationService interface {
	CreateOrganization(ctx context.Context, req *models.CreateOrganizationRequest) (*models.Organization, error)
	GetOrganization(ctx context.Context, id string) (*models.Organization, error)
	ListOrganizations(ctx context.Context) ([]models.Organization, error)
}

// organizationService is the concrete implementation of OrganizationService
type organizationService struct {
	orgRepo repositories.OrganizationRepository
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(orgRepo repositories.OrganizationRepository) OrganizationService {
	return &organizationService{
		orgRepo: orgRepo,
	}
}

// CreateOrganization creates a new organization
func (s *organizationService) CreateOrganization(ctx context.Context, req *models.CreateOrganizationRequest) (*models.Organization, error) {
	org := &models.Organization{
		Name: req.Name,
	}

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganization retrieves an organization by ID
func (s *organizationService) GetOrganization(ctx context.Context, id string) (*models.Organization, error) {
	parsedUUID, _ := uuid.Parse(id)
	return s.orgRepo.FindByID(ctx, parsedUUID)
}

// ListOrganizations retrieves all organizations
func (s *organizationService) ListOrganizations(ctx context.Context) ([]models.Organization, error) {
	return s.orgRepo.FindAll(ctx)
}

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
	GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	ListOrganizations(ctx context.Context) ([]models.Organization, error)
	UpdateOrganization(ctx context.Context, id uuid.UUID, req *models.UpdateOrganizationRequest) (*models.Organization, error)
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
	org := req.ToDomain()

	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	return s.orgRepo.FindByID(ctx, org.ID)
}

// GetOrganization retrieves an organization by ID
func (s *organizationService) GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return s.orgRepo.FindByID(ctx, id)
}

// ListOrganizations retrieves all organizations
func (s *organizationService) ListOrganizations(ctx context.Context) ([]models.Organization, error) {
	return s.orgRepo.FindAll(ctx)
}

func (s *organizationService) UpdateOrganization(ctx context.Context, id uuid.UUID, req *models.UpdateOrganizationRequest) (*models.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		org.Name = *req.Name
	}
	if req.Latitude != nil && req.Longitude != nil {
		org.Geom.Latitude = *req.Latitude
		org.Geom.Longitude = *req.Longitude
	}
	if req.Address != nil {
		org.Address = req.Address
	}
	if req.Description != nil {
		org.Description = req.Description
	}
	if req.ImageUrls != nil {
		org.ImageUrls = req.ImageUrls
	}

	return org, s.orgRepo.Update(ctx, org)
}

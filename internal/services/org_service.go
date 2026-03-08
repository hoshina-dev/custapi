package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"github.com/hoshina-dev/custapi/internal/repositories"
)

// OrganizationService defines organization business logic operations
type OrganizationService interface {
	CreateOrganization(ctx context.Context, req *models.CreateOrganizationRequest) (*models.Organization, error)
	GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	GetByIDs(ctx context.Context, id []uuid.UUID) ([]models.Organization, error)
	ListOrganizations(ctx context.Context) ([]models.Organization, error)
	GetAllCoords(ctx context.Context) ([]models.Organization, error)
	UpdateOrganization(ctx context.Context, id uuid.UUID, req *models.UpdateOrganizationRequest) (*models.Organization, error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	SearchOrganizations(ctx context.Context, query string, limit int) ([]models.Organization, error)
	ListMembers(ctx context.Context, orgID uuid.UUID) ([]models.UserOrganization, error)
	AddMember(ctx context.Context, orgID uuid.UUID, req *models.AddMemberRequest) error
	RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error
	SetRole(ctx context.Context, orgID, userID uuid.UUID, role models.MemberRole) error
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

	return org, nil
}

// GetOrganization retrieves an organization by ID
func (s *organizationService) GetOrganization(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	return s.orgRepo.FindByID(ctx, id)
}

func (s *organizationService) GetByIDs(ctx context.Context, id []uuid.UUID) ([]models.Organization, error) {
	return s.orgRepo.FindByIDs(ctx, id)
}

// ListOrganizations retrieves all organizations
func (s *organizationService) ListOrganizations(ctx context.Context) ([]models.Organization, error) {
	return s.orgRepo.FindAll(ctx)
}

func (s *organizationService) GetAllCoords(ctx context.Context) ([]models.Organization, error) {
	return s.orgRepo.FindAllCoords(ctx)
}

func (s *organizationService) UpdateOrganization(ctx context.Context, id uuid.UUID, req *models.UpdateOrganizationRequest) (*models.Organization, error) {
	org, err := s.orgRepo.FindByID(ctx, id)
	if err != nil || org == nil {
		return nil, err
	}

	updatedOrg := req.ToDomain(org.ID)

	return updatedOrg, s.orgRepo.Update(ctx, updatedOrg)
}

func (s *organizationService) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	return s.orgRepo.Delete(ctx, id)
}

// SearchOrganizations searches organizations by name
func (s *organizationService) SearchOrganizations(ctx context.Context, query string, limit int) ([]models.Organization, error) {
	return s.orgRepo.Search(ctx, query, limit)
}

// ListMembers returns all members of an organization
func (s *organizationService) ListMembers(ctx context.Context, orgID uuid.UUID) ([]models.UserOrganization, error) {
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New("organization not found")
	}
	return s.orgRepo.FindMembers(ctx, orgID)
}

// AddMember adds a user to an organization
func (s *organizationService) AddMember(ctx context.Context, orgID uuid.UUID, req *models.AddMemberRequest) error {
	org, err := s.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return err
	}
	if org == nil {
		return errors.New("organization not found")
	}
	return s.orgRepo.AddMember(ctx, &models.UserOrganization{
		OrganizationID: orgID,
		UserID:         req.UserID,
		Role:           req.Role,
	})
}

// RemoveMember removes a user from an organization
func (s *organizationService) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	return s.orgRepo.RemoveMember(ctx, orgID, userID)
}

// SetRole updates the role of a member
func (s *organizationService) SetRole(ctx context.Context, orgID, userID uuid.UUID, role models.MemberRole) error {
	return s.orgRepo.SetRole(ctx, orgID, userID, role)
}

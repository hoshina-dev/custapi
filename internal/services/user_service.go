package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"github.com/hoshina-dev/custapi/internal/repositories"
)

// UserService defines user business logic operations
type UserService interface {
	CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	ListUsersByOrganization(ctx context.Context, orgID string) ([]models.User, error)
	SearchUsers(ctx context.Context, query string) ([]models.User, error)
}

// userService is the concrete implementation of UserService
type userService struct {
	userRepo repositories.UserRepository
	orgRepo  repositories.OrganizationRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository, orgRepo repositories.OrganizationRepository) UserService {
	return &userService{
		userRepo: userRepo,
		orgRepo:  orgRepo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Verify organization exists
	parsedUUID, _ := uuid.Parse(req.OrganizationID)
	org, err := s.orgRepo.FindByID(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New("organization not found")
	}

	user := &models.User{
		Email:          req.Email,
		Name:           req.Name,
		OrganizationID: req.OrganizationID,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, id string) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// ListUsers retrieves all users
func (s *userService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

// ListUsersByOrganization retrieves users by organization
func (s *userService) ListUsersByOrganization(ctx context.Context, orgID string) ([]models.User, error) {
	// Verify organization exists
	parsedUUID, _ := uuid.Parse(orgID)
	org, err := s.orgRepo.FindByID(ctx, parsedUUID)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, errors.New("organization not found")
	}

	return s.userRepo.FindByOrganizationID(ctx, orgID)
}

// SearchUsers searches users by name or email
func (s *userService) SearchUsers(ctx context.Context, query string) ([]models.User, error) {
	return s.userRepo.Search(ctx, query)
}

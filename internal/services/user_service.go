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
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]models.UserOrganization, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error)
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
	user, err := req.ToDomain()
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetUserByEmail retrieves a user by email address
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// ListUsers retrieves all users
func (s *userService) ListUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

// GetUserOrganizations retrieves all organizations a user belongs to
func (s *userService) GetUserOrganizations(ctx context.Context, userID uuid.UUID) ([]models.UserOrganization, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return s.userRepo.FindUserOrganizations(ctx, userID)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil || user == nil {
		return nil, err
	}

	updatedUser, err := req.ToDomain(user.ID)
	if err != nil {
		return nil, err
	}

	return updatedUser, s.userRepo.Update(ctx, updatedUser)
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}

// SearchUsers searches users by name or email
func (s *userService) SearchUsers(ctx context.Context, query string, limit int) ([]models.User, error) {
	return s.userRepo.Search(ctx, query, limit)
}

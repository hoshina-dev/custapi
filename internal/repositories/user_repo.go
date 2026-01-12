package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"gorm.io/gorm"
)

// UserRepository defines user persistence operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id string) (*models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	FindByOrganizationID(ctx context.Context, orgID string) ([]models.User, error)
}

// userRepository is the concrete implementation of UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.New().String()
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Organization").First(&user, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll retrieves all users
func (r *userRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Preload("Organization").Order("created_at DESC").Find(&users).Error
	return users, err
}

// FindByOrganizationID finds all users in an organization
func (r *userRepository) FindByOrganizationID(ctx context.Context, orgID string) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Order("created_at DESC").Find(&users).Error
	return users, err
}

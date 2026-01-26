package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRepository defines user persistence operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindAll(ctx context.Context) ([]models.User, error)
	FindByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string) ([]models.User, error)
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
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
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
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&users).Error
	return users, err
}

// FindByOrganizationID finds all users in an organization
func (r *userRepository) FindByOrganizationID(ctx context.Context, orgID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).Where("organization_id = ?", orgID).Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Model(user).Clauses(clause.Returning{}).Updates(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

// Search searches users by name or email using ILIKE
func (r *userRepository) Search(ctx context.Context, query string, limit int) ([]models.User, error) {
	var users []models.User
	searchPattern := "%" + query + "%"
	db := r.db.WithContext(ctx).
		Preload("Organization").
		Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern).
		Order("name ASC")

	if limit > 0 {
		db = db.Limit(limit)
	}

	err := db.Find(&users).Error
	return users, err
}

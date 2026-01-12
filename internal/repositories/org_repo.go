package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"gorm.io/gorm"
)

// OrganizationRepository defines organization persistence operations
type OrganizationRepository interface {
	Create(ctx context.Context, org *models.Organization) error
	FindByID(ctx context.Context, id string) (*models.Organization, error)
	FindAll(ctx context.Context) ([]models.Organization, error)
}

// organizationRepository is the concrete implementation of OrganizationRepository
type organizationRepository struct {
	db *gorm.DB
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

// Create creates a new organization
func (r *organizationRepository) Create(ctx context.Context, org *models.Organization) error {
	org.ID = uuid.New().String()
	return r.db.WithContext(ctx).Create(org).Error
}

// FindByID finds an organization by ID
func (r *organizationRepository) FindByID(ctx context.Context, id string) (*models.Organization, error) {
	var org models.Organization
	err := r.db.WithContext(ctx).First(&org, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// FindAll retrieves all organizations
func (r *organizationRepository) FindAll(ctx context.Context) ([]models.Organization, error) {
	var orgs []models.Organization
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&orgs).Error
	return orgs, err
}

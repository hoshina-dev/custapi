package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoshina-dev/custapi/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrganizationRepository defines organization persistence operations
type OrganizationRepository interface {
	Create(ctx context.Context, org *models.Organization) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Organization, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Organization, error)
	FindAll(ctx context.Context) ([]models.Organization, error)
	FindAllCoords(ctx context.Context) ([]models.Organization, error)
	Update(ctx context.Context, org *models.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string, limit int) ([]models.Organization, error)
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
	return r.db.WithContext(ctx).Create(org).Error
}

// FindByID finds an organization by ID
func (r *organizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Organization, error) {
	var org models.Organization
	err := r.db.WithContext(ctx).First(&org, id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.Organization, error) {
	var orgs []models.Organization
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Order("created_at DESC").Find(&orgs).Error
	return orgs, err
}

// FindAll retrieves all organizations
func (r *organizationRepository) FindAll(ctx context.Context) ([]models.Organization, error) {
	var orgs []models.Organization
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&orgs).Error
	return orgs, err
}

func (r *organizationRepository) FindAllCoords(ctx context.Context) ([]models.Organization, error) {
	var orgs []models.Organization
	err := r.db.WithContext(ctx).Select("id, latitude, longitude").Find(&orgs).Error
	return orgs, err
}

func (r *organizationRepository) Update(ctx context.Context, org *models.Organization) error {
	return r.db.WithContext(ctx).Model(org).Clauses(clause.Returning{}).Updates(org).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&models.Organization{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("organization not found")
	}
	return nil
}

// Search searches organizations by name using ILIKE
func (r *organizationRepository) Search(ctx context.Context, query string, limit int) ([]models.Organization, error) {
	var orgs []models.Organization
	searchPattern := "%" + query + "%"
	db := r.db.WithContext(ctx).
		Where("name ILIKE ?", searchPattern).
		Order("name ASC")

	if limit > 0 {
		db = db.Limit(limit)
	}

	err := db.Find(&orgs).Error
	return orgs, err
}

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
	FindMembers(ctx context.Context, orgID uuid.UUID) ([]models.User, error)
	AddMember(ctx context.Context, m *models.UserOrganization) error
	RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error
	SetRole(ctx context.Context, orgID, userID uuid.UUID, role models.MemberRole) error
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
	escaped := escapeLike(query)
	searchPattern := "%" + escaped + "%"
	db := r.db.WithContext(ctx).
		Where("name ILIKE ? ESCAPE '\\'", searchPattern).
		Order("name ASC")

	if limit > 0 {
		db = db.Limit(limit)
	}

	err := db.Find(&orgs).Error
	return orgs, err
}

// FindMembers returns all non-deleted users belonging to an organization
func (r *organizationRepository) FindMembers(ctx context.Context, orgID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_organizations ON user_organizations.user_id = users.id").
		Where("user_organizations.organization_id = ?", orgID).
		Preload("Organizations.Organization").
		Order("users.name ASC").
		Find(&users).Error
	return users, err
}

// AddMember inserts a row into user_organizations
func (r *organizationRepository) AddMember(ctx context.Context, m *models.UserOrganization) error {
	return r.db.WithContext(ctx).Create(m).Error
}

// RemoveMember deletes a user_organizations row (hard delete — no deleted_at on junction table)
func (r *organizationRepository) RemoveMember(ctx context.Context, orgID, userID uuid.UUID) error {
	res := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Delete(&models.UserOrganization{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return nil
}

// SetRole updates the role for an existing membership
func (r *organizationRepository) SetRole(ctx context.Context, orgID, userID uuid.UUID, role models.MemberRole) error {
	res := r.db.WithContext(ctx).
		Model(&models.UserOrganization{}).
		Where("organization_id = ? AND user_id = ?", orgID, userID).
		Update("role", role)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("member not found")
	}
	return nil
}

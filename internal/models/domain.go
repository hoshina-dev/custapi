package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// MemberRole represents a user's role within an organization
type MemberRole string // @name MemberRole

const (
	RoleAdmin MemberRole = "admin"
	RoleUser  MemberRole = "user"
)

// User represents a user in the system
type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email              string
	Name               string
	Password           string
	PhoneNumber        *string
	SocialMedia        *string
	Description        *string
	AvatarURL          *string
	ResearchCategories pq.StringArray     `gorm:"type:text[];default:'{}'"`
	Organizations      []UserOrganization `gorm:"foreignKey:UserID"`
	CreatedAt          time.Time          `gorm:"autoCreateTime"`
	UpdatedAt          time.Time          `gorm:"autoUpdateTime"`
	DeletedAt          gorm.DeletedAt
}

// Organization represents an organization in the system
type Organization struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string
	Latitude    *float64
	Longitude   *float64
	Address     *string
	Description *string
	ImageUrls   pq.StringArray `gorm:"type:text[];default:'{}'"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
}

// UserOrganization represents a user's membership in an organization.
type UserOrganization struct {
	UserID         uuid.UUID    `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Role           MemberRole   `gorm:"type:varchar(50);default:'user'"`
	CreatedAt      time.Time    `gorm:"autoCreateTime"`
	User           User         `gorm:"foreignKey:UserID"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
}

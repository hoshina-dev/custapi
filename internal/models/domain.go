package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email              string
	Name               string
	OrganizationID     uuid.UUID
	Organization       Organization
	Password           string
	IsAdmin            bool
	PhoneNumber        *string
	SocialMedia        *string
	Description        *string
	AvatarURL          *string
	ResearchCategories pq.StringArray
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
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
	Users       []User
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
}

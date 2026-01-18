package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// User represents a user in the system
type User struct {
	ID             string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email          string         `gorm:"uniqueIndex;not null" json:"email"`
	Name           string         `gorm:"not null" json:"name"`
	OrganizationID string         `gorm:"type:uuid;not null;index" json:"organization_id"`
	Organization   *Organization  `gorm:"foreignKey:OrganizationID" json:"organization,omitempty"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// Organization represents an organization in the system
type Organization struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string
	Geom        Point
	Address     *string
	Description *string
	ImageUrls   pq.StringArray `gorm:"type:text[];default:'{}'"`
	Users       []User
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt
}

type Point struct {
	Latitude  *float64
	Longitude *float64
	Geom      string
}

func (p *Point) Scan(v interface{}) error {
	if v != nil {
		p.Geom = v.(string)
	}
	return nil
}

func (p Point) GormDataType() string {
	return "geometry(Point,4326)"
}

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return gorm.Expr("ST_Point(?, ?, 4326)", p.Latitude, p.Longitude)
}

func (o *Organization) AfterSave(tx *gorm.DB) (err error) {
	err = tx.First(o, o.ID).Error
	return
}

func (o *Organization) AfterFind(tx *gorm.DB) (err error) {
	if o.Geom.Geom != "" {
		var tmp string
		err = tx.Raw("SELECT ST_AsGeoJSON(?)", o.Geom.Geom).Row().Scan(&tmp)
		if err == nil {
			o.Geom.Geom = tmp
		}
	}
	return
}

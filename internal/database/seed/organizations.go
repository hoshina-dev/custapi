package seed

import (
	"github.com/hoshina-dev/custapi/internal/models"
	"gorm.io/gorm"
)

func seedOrganizations(db *gorm.DB) error {
	orgs := []models.Organization{
		{
			Name:        "Chulalongkorn University",
			Description: strPtr("Thailand's leading university, dedicated to academic excellence, research innovation, and societal impact."),
			Latitude:    floatPtr(13.73826),
			Longitude:   floatPtr(100.532413),
		},
		{
			Name:        "The National Institutes for Quantum Science and Technology (QST)",
			Description: strPtr("Japan's leading institute for quantum science and technology research and development."),
			Latitude:    floatPtr(35.6350881070613),
			Longitude:   floatPtr(140.1024440117597),
		},
	}

	for _, org := range orgs {
		if err := db.Where("name = ?", org.Name).FirstOrCreate(&org).Error; err != nil {
			return err
		}
	}

	return nil
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}

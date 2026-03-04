package seed

import (
	"github.com/hoshina-dev/custapi/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedUsers(db *gorm.DB) error {
	var orgs []models.Organization
	if err := db.Find(&orgs).Error; err != nil {
		return err
	}

	type userSeed struct {
		user    models.User
		orgIdx  int
		isAdmin bool
	}

	seeds := []userSeed{
		{user: models.User{Email: "admin@chula.test", Name: "Chula Admin", Password: hash("password123")}, orgIdx: 0, isAdmin: true},
		{user: models.User{Email: "user1@chula.test", Name: "Researcher 1", Password: hash("321drowssap")}, orgIdx: 0, isAdmin: false},
		{user: models.User{Email: "researcher@qst.test", Name: "Sam the Scientist", Password: hash("123password")}, orgIdx: 1, isAdmin: false},
	}

	for i := range seeds {
		u := &seeds[i].user
		if err := db.Where("email = ?", u.Email).FirstOrCreate(u).Error; err != nil {
			return err
		}
		membership := models.UserOrganization{
			UserID:         u.ID,
			OrganizationID: orgs[seeds[i].orgIdx].ID,
			IsAdmin:        seeds[i].isAdmin,
		}
		if err := db.Where("user_id = ? AND organization_id = ?", membership.UserID, membership.OrganizationID).
			FirstOrCreate(&membership).Error; err != nil {
			return err
		}
	}

	return nil
}

func hash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "examplehashedpassword"
	}
	return string(hashedPassword)
}

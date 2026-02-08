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

	users := []models.User{
		{
			Email:          "admin@chula.test",
			Name:           "Chula Admin",
			IsAdmin:        true,
			OrganizationID: orgs[0].ID,
			Password:       hash("password123"),
		},
		{
			Email:          "user1@chula.test",
			Name:           "Researcher 1",
			IsAdmin:        false,
			OrganizationID: orgs[0].ID,
			Password:       hash("321drowssap"),
		},
		{
			Email:          "researcher@qst.test",
			Name:           "Sam the Scientist",
			IsAdmin:        false,
			OrganizationID: orgs[1].ID,
			Password:       hash("123password"),
		},
	}

	for _, user := range users {
		if err := db.Where("email = ?", user.Email).FirstOrCreate(&user).Error; err != nil {
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

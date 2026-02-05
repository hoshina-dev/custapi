package seed

import "gorm.io/gorm"

func Run(db *gorm.DB) error {
	if err := seedOrganizations(db); err != nil {
		return err
	}

	if err := seedUsers(db); err != nil {
		return err
	}

	return nil
}

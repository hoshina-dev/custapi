package main

import (
	"log"

	"github.com/hoshina-dev/custapi/internal/config"
	"github.com/hoshina-dev/custapi/internal/database"
	"github.com/hoshina-dev/custapi/internal/database/seed"
)

func main() {
	cfg := config.Load()

	db := database.ConnectDB(cfg.DataSourceName)

	if err := seed.Run(db); err != nil {
		log.Fatal(err)
	}

	log.Println("Dev database seeded successfully")
}

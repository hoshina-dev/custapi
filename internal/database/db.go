package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

func ConnectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Emits a child span (parented to whatever span is on the query's
	// context.Context, e.g. the HTTP request span from otelfiber) plus
	// db.client.* metrics for every Create/Query/Update/Delete/Row/Raw call.
	// Callers must use db.WithContext(ctx) for spans to be parented correctly.
	if err := db.Use(tracing.NewPlugin()); err != nil {
		log.Fatalf("Failed to register OpenTelemetry GORM plugin: %v", err)
	}

	return db
}

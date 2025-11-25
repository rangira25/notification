package storage

import (
	"log"
	"time"

	"github.com/rangira25/notification/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGorm(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed connect db: %v", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS citext;")


	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(5 * time.Minute)
	}

	// AutoMigrate for dev. Replace with migrations in prod.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("auto migrate err: %v", err)
	}
	return db
}

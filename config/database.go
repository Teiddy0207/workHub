package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDatabase trả về *gorm.DB và error
func ConnectDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_MASTER_HOST"),
		os.Getenv("DB_MASTER_USER"),
		os.Getenv("DB_MASTER_PASSWORD"),
		os.Getenv("DB_MASTER_NAME"),
		os.Getenv("DB_MASTER_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	return db, nil
}

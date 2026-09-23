// Package storage manages the SQLite database persistence and schema migrations.
package storage

import (
	"fmt"

	"facturacion-local/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the SQLite database connection and runs schema auto-migrations.
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Run auto-migrations for domain models
	err = db.AutoMigrate(&models.Client{}, &models.Invoice{}, &models.InvoiceLine{})
	if err != nil {
		return nil, fmt.Errorf("failed to run database schema migrations: %w", err)
	}

	return db, nil
}

package storage

import (
	"facturacion-local/internal/models"
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// InitDB inicializa la conexión a la base de datos SQLite y ejecuta las migraciones.
func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error al conectar con SQLite: %w", err)
	}

	// Ejecutar AutoMigrate para generar/actualizar el esquema
	err = db.AutoMigrate(&models.Client{}, &models.Invoice{}, &models.InvoiceLine{})
	if err != nil {
		return nil, fmt.Errorf("error al migrar la base de datos: %w", err)
	}

	return db, nil
}

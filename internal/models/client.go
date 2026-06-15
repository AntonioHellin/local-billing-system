package models

import (
	"time"

	"gorm.io/gorm"
)

// Client representa un cliente dentro del sistema de facturación.
type Client struct {
	ID        uint           `gorm:"primarykey"`
	Name      string         `gorm:"not null"`
	TaxID     string         `gorm:"not null;uniqueIndex"` // NIF/CIF
	Address   string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

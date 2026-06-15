package models

import (
	"time"

	"gorm.io/gorm"
)

// InvoiceLine representa una línea de detalle dentro de una factura.
type InvoiceLine struct {
	ID          uint           `gorm:"primarykey"`
	InvoiceID   uint           `gorm:"not null"`
	Description string         `gorm:"not null"`
	Quantity    float64        `gorm:"not null"`
	UnitPrice   float64        `gorm:"not null"`
	TaxRate     float64        `gorm:"not null"`
	Total       float64        `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

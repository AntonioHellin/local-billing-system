package models

import (
	"time"

	"gorm.io/gorm"
)

// Invoice representa una factura en el sistema.
// La combinación de InvoiceNumber y Version debe ser única para soportar facturas rectificativas.
type Invoice struct {
	ID            uint           `gorm:"primarykey"`
	ClientID      uint           `gorm:"not null"`
	Client        Client         `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	InvoiceNumber string         `gorm:"uniqueIndex:idx_invoice_number_version,not null"`
	Version       int            `gorm:"uniqueIndex:idx_invoice_number_version,not null;default:1"`
	IssueDate     time.Time      `gorm:"not null"`
	DueDate       time.Time
	Subtotal      float64        `gorm:"not null"`
	TaxTotal      float64        `gorm:"not null"`
	Total         float64        `gorm:"not null"`
	Status        string         `gorm:"not null;default:'DRAFT'"`
	Notes         string
	Lines         []InvoiceLine  `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

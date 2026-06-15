package storage_test

import (
	"facturacion-local/internal/models"
	"facturacion-local/internal/storage"
	"testing"
	"time"
)

func TestInitDBAndConstraints(t *testing.T) {
	// Usamos una base de datos en memoria para los tests
	db, err := storage.InitDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("Fallo al inicializar la BD en memoria: %v", err)
	}

	// 1. Crear un Cliente
	client := models.Client{
		Name:  "Test Client",
		TaxID: "12345678A",
	}
	if err := db.Create(&client).Error; err != nil {
		t.Fatalf("Fallo al crear cliente: %v", err)
	}

	// 2. Crear una Factura
	invoice1 := models.Invoice{
		ClientID:      client.ID,
		InvoiceNumber: "INV-001",
		Version:       1,
		IssueDate:     time.Now(),
		Subtotal:      100.0,
		TaxTotal:      21.0,
		Total:         121.0,
	}
	if err := db.Create(&invoice1).Error; err != nil {
		t.Fatalf("Fallo al crear factura 1: %v", err)
	}

	// 3. Crear otra factura con el mismo número pero distinta versión (Rectificativa)
	invoice2 := models.Invoice{
		ClientID:      client.ID,
		InvoiceNumber: "INV-001",
		Version:       2,
		IssueDate:     time.Now(),
		Subtotal:      -100.0,
		TaxTotal:      -21.0,
		Total:         -121.0,
	}
	if err := db.Create(&invoice2).Error; err != nil {
		t.Fatalf("Fallo al crear factura 2 (versión 2): %v", err)
	}

	// 4. Intentar crear una factura con el MISMO número y MISMA versión (Debe fallar)
	invoice3 := models.Invoice{
		ClientID:      client.ID,
		InvoiceNumber: "INV-001",
		Version:       1,
		IssueDate:     time.Now(),
		Subtotal:      50.0,
		TaxTotal:      10.5,
		Total:         60.5,
	}
	err = db.Create(&invoice3).Error
	if err == nil {
		t.Fatalf("Se esperaba un error al violar la restricción de unicidad (InvoiceNumber + Version), pero no ocurrió")
	}
}

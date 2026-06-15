package main

import (
	"fmt"
	"time"

	"facturacion-local/internal/exporter"
	"facturacion-local/internal/models"
)

func main() {
	// Crear cliente de prueba
	client := models.Client{
		Name:    "Empresa Cliente S.A.",
		TaxID:   "A98765432",
		Address: "Av. Principal 456, Barcelona",
		Email:   "contacto@empresacliente.es",
	}

	// Crear líneas de factura
	lines := []models.InvoiceLine{
		{
			Description: "Consultoría IT",
			Quantity:    10,
			UnitPrice:   60.00,
			TaxRate:     21,
			Total:       600.00,
		},
		{
			Description: "Licencia Anual Software",
			Quantity:    1,
			UnitPrice:   1200.00,
			TaxRate:     21,
			Total:       1200.00,
		},
	}

	// Calcular totales
	subtotal := 0.0
	for _, line := range lines {
		subtotal += line.Total
	}
	taxTotal := subtotal * 0.21
	total := subtotal + taxTotal

	// Crear factura consolidada
	invoice := &models.Invoice{
		InvoiceNumber: "INV-2026-001",
		Version:       1,
		IssueDate:     time.Now(),
		Client:        client,
		Lines:         lines,
		Subtotal:      subtotal,
		TaxTotal:      taxTotal,
		Total:         total,
		Notes:         "Pago a 30 días mediante transferencia bancaria.",
	}

	// Inicializar el exportador y generar PDF
	exp := exporter.NewPDFExporter("") // Usará ruta por defecto
	path, err := exp.ExportInvoice(invoice)
	if err != nil {
		fmt.Printf("Error al exportar la factura: %v\n", err)
		return
	}

	fmt.Printf("Factura exportada con éxito en: %s\n", path)
}

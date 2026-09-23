// Package exporter handles generation and export of invoice PDF documents.
package exporter

import (
	"fmt"
	"os"
	"path/filepath"

	"facturacion-local/internal/models"

	"github.com/jung-kurt/gofpdf"
)

// DefaultExportPath returns the directory path for exported invoices.
// It prioritizes the EXPORT_PATH environment variable, falling back to a subfolder in the user's home directory.
func DefaultExportPath() string {
	if envPath := os.Getenv("EXPORT_PATH"); envPath != "" {
		return envPath
	}
	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, "BillingExports")
	}
	return filepath.Join(".", "exports")
}

// PDFExporter manages generation and formatting of PDF invoices.
type PDFExporter struct {
	ExportPath string
}

// NewPDFExporter creates a new PDFExporter instance.
// If exportPath is empty, DefaultExportPath() is used.
func NewPDFExporter(exportPath string) *PDFExporter {
	if exportPath == "" {
		exportPath = DefaultExportPath()
	}
	return &PDFExporter{
		ExportPath: exportPath,
	}
}

// ExportInvoice renders a consolidated invoice into a formatted PDF file.
func (e *PDFExporter) ExportInvoice(invoice *models.Invoice) (string, error) {
	// Verify directory exists or create it
	if err := os.MkdirAll(e.ExportPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create or verify export directory: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	// --- Header (Issuer) ---
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Local Billing Services")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Tax ID: B12345678")
	pdf.Ln(5)
	pdf.Cell(40, 10, "Main St 123, Madrid")
	pdf.Ln(15)

	// --- Invoice Title ---
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fmt.Sprintf("INVOICE: %s (v%d)", invoice.InvoiceNumber, invoice.Version))
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Issue Date: %s", invoice.IssueDate.Format("02-01-2006")))
	pdf.Ln(12)

	// --- Client Details ---
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Client Details:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Name: %s", invoice.Client.Name))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Tax ID: %s", invoice.Client.TaxID))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Address: %s", invoice.Client.Address))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Email: %s", invoice.Client.Email))
	pdf.Ln(15)

	// --- Line Items Table ---
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(80, 10, "Description", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 10, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Unit Price", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 10, "Tax %", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Total", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 10)
	for _, line := range invoice.Lines {
		pdf.CellFormat(80, 10, line.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 10, fmt.Sprintf("%.2f", line.Quantity), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f EUR", line.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(20, 10, fmt.Sprintf("%.0f%%", line.TaxRate), "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", line.Total), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(10)

	// --- Totals ---
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(110, 10, "")
	pdf.CellFormat(40, 10, "Subtotal:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", invoice.Subtotal), "1", 1, "R", false, 0, "")

	pdf.Cell(110, 10, "")
	pdf.CellFormat(40, 10, "Tax Total:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", invoice.TaxTotal), "1", 1, "R", false, 0, "")

	pdf.Cell(110, 10, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "TOTAL:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", invoice.Total), "1", 1, "R", false, 0, "")

	pdf.Ln(10)
	if invoice.Notes != "" {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, "Notes: "+invoice.Notes, "", "L", false)
	}

	// --- File Output ---
	filename := fmt.Sprintf("Invoice_%s_v%d.pdf", invoice.InvoiceNumber, invoice.Version)
	fullPath := filepath.Join(e.ExportPath, filename)

	if err := pdf.OutputFileAndClose(fullPath); err != nil {
		return "", fmt.Errorf("failed to save PDF invoice: %w", err)
	}

	return fullPath, nil
}

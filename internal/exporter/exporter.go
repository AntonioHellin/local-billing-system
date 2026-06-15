package exporter

import (
	"fmt"
	"os"
	"path/filepath"

	"facturacion-local/internal/models"

	"github.com/jung-kurt/gofpdf"
)

const defaultExportPath = `C:\Facturacion_Local\Exportaciones`

// PDFExporter maneja la generación de documentos PDF.
type PDFExporter struct {
	ExportPath string
}

// NewPDFExporter crea una nueva instancia de PDFExporter.
// Si no se proporciona un path, se utiliza el defaultExportPath.
func NewPDFExporter(exportPath string) *PDFExporter {
	if exportPath == "" {
		exportPath = defaultExportPath
	}
	return &PDFExporter{
		ExportPath: exportPath,
	}
}

// ExportInvoice toma la estructura de una factura consolidada y la exporta a PDF.
func (e *PDFExporter) ExportInvoice(invoice *models.Invoice) (string, error) {
	// Verificar si la carpeta existe o crearla
	if err := os.MkdirAll(e.ExportPath, 0755); err != nil {
		return "", fmt.Errorf("error al crear o verificar directorio destino: %w", err)
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)

	// --- Header (Emisor) ---
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Mi Empresa S.L.")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "NIF: B12345678")
	pdf.Ln(5)
	pdf.Cell(40, 10, "Calle Falsa 123, Madrid")
	pdf.Ln(15)

	// --- Título de Factura ---
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, fmt.Sprintf("FACTURA: %s (v%d)", invoice.InvoiceNumber, invoice.Version))
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Fecha Emision: %s", invoice.IssueDate.Format("02-01-2006")))
	pdf.Ln(12)

	// --- Datos del Cliente ---
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Datos del Cliente:")
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 10, fmt.Sprintf("Nombre: %s", invoice.Client.Name))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("NIF/CIF: %s", invoice.Client.TaxID))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Direccion: %s", invoice.Client.Address))
	pdf.Ln(5)
	pdf.Cell(0, 10, fmt.Sprintf("Email: %s", invoice.Client.Email))
	pdf.Ln(15)

	// --- Tabla de Conceptos ---
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(80, 10, "Descripcion", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 10, "Cant.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Precio U.", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 10, "% IVA", "1", 0, "C", false, 0, "")
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

	// --- Totales ---
	pdf.SetFont("Arial", "B", 10)
	// Offset para alinear a la derecha
	pdf.Cell(110, 10, "")
	pdf.CellFormat(40, 10, "Subtotal:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", invoice.Subtotal), "1", 1, "R", false, 0, "")

	pdf.Cell(110, 10, "")
	pdf.CellFormat(40, 10, "Total Impuestos:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", invoice.TaxTotal), "1", 1, "R", false, 0, "")

	pdf.Cell(110, 10, "")
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 10, "TOTAL:", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f EUR", invoice.Total), "1", 1, "R", false, 0, "")

	pdf.Ln(10)
	if invoice.Notes != "" {
		pdf.SetFont("Arial", "I", 9)
		pdf.MultiCell(0, 5, "Notas: "+invoice.Notes, "", "L", false)
	}

	// --- Guardado ---
	filename := fmt.Sprintf("Factura_%s_v%d.pdf", invoice.InvoiceNumber, invoice.Version)
	fullPath := filepath.Join(e.ExportPath, filename)

	if err := pdf.OutputFileAndClose(fullPath); err != nil {
		return "", fmt.Errorf("error al guardar PDF: %w", err)
	}

	return fullPath, nil
}

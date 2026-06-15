package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"facturacion-local/internal/models"
)

// CreateInvoiceRequest represents the expected payload for POST /api/invoices
type CreateInvoiceRequest struct {
	ClientID      uint    `json:"client_id"`
	BaseImponible float64 `json:"base_imponible"`
	InvoiceNumber string  `json:"invoice_number,omitempty"`
}

// CreateInvoice handles POST /api/invoices
func (a *API) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var req CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Strict Validations
	if req.ClientID <= 0 {
		respondError(w, http.StatusBadRequest, "client_id is required and must be greater than 0")
		return
	}
	if req.BaseImponible <= 0 {
		respondError(w, http.StatusBadRequest, "base_imponible is required and must be greater than 0")
		return
	}

	// Check if client exists to avoid FK constraint errors
	var client models.Client
	if err := a.DB.First(&client, req.ClientID).Error; err != nil {
		respondError(w, http.StatusNotFound, "Client not found")
		return
	}

	// Tax Calculation logic
	subtotal := req.BaseImponible
	taxRate := 0.21 // 21% IVA
	taxTotal := subtotal * taxRate
	total := subtotal + taxTotal

	// Generate a temporary InvoiceNumber if not provided
	invoiceNum := req.InvoiceNumber
	if invoiceNum == "" {
		invoiceNum = fmt.Sprintf("INV-%d", time.Now().Unix())
	}

	invoice := models.Invoice{
		ClientID:      req.ClientID,
		InvoiceNumber: invoiceNum,
		Version:       1,
		IssueDate:     time.Now(),
		Subtotal:      subtotal,
		TaxTotal:      taxTotal,
		Total:         total,
		Status:        "DRAFT",
	}

	if err := a.DB.Create(&invoice).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating invoice")
		return
	}

	// Preload client for the response so the frontend has complete info
	a.DB.Preload("Client").First(&invoice, invoice.ID)

	respondJSON(w, http.StatusCreated, invoice)
}

// GetClientInvoices handles GET /api/clients/{id}/invoices
func (a *API) GetClientInvoices(w http.ResponseWriter, r *http.Request) {
	clientID := r.PathValue("id")
	if clientID == "" {
		respondError(w, http.StatusBadRequest, "Missing client ID")
		return
	}

	var invoices []models.Invoice
	if err := a.DB.Preload("Client").Where("client_id = ?", clientID).Order("issue_date desc").Find(&invoices).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving invoices")
		return
	}

	if invoices == nil {
		invoices = []models.Invoice{}
	}

	respondJSON(w, http.StatusOK, invoices)
}

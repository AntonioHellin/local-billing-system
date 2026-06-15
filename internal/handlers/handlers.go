package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

// API encapsulates the database connection and handlers
type API struct {
	DB *gorm.DB
}

// RegisterRoutes registers the HTTP endpoints on the provided mux
func RegisterRoutes(mux *http.ServeMux, db *gorm.DB) {
	api := &API{DB: db}

	// GET /api/clients
	mux.HandleFunc("GET /api/clients", api.GetClients)

	// POST /api/invoices
	mux.HandleFunc("POST /api/invoices", api.CreateInvoice)

	// POST /api/clients
	mux.HandleFunc("POST /api/clients", api.CreateClient)

	// GET /api/clients/{id}/invoices
	mux.HandleFunc("GET /api/clients/{id}/invoices", api.GetClientInvoices)
}

// respondJSON is a helper to send JSON responses
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
	}
}

// respondError is a helper to send error messages as JSON
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

package handlers

import (
	"encoding/json"
	"net/http"

	"facturacion-local/internal/models"
)

// GetClients handles GET /api/clients
func (a *API) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []models.Client
	
	if err := a.DB.Find(&clients).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving clients")
		return
	}

	// If no clients, clients will be an empty slice (represented as [] in JSON instead of null)
	if clients == nil {
		clients = []models.Client{}
	}

	respondJSON(w, http.StatusOK, clients)
}

// CreateClientRequest represents the expected payload for POST /api/clients
type CreateClientRequest struct {
	Name    string `json:"name"`
	TaxID   string `json:"tax_id"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// CreateClient handles POST /api/clients
func (a *API) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req CreateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	if req.Name == "" || req.TaxID == "" {
		respondError(w, http.StatusBadRequest, "Name and TaxID are required")
		return
	}

	client := models.Client{
		Name:    req.Name,
		TaxID:   req.TaxID,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}

	if err := a.DB.Create(&client).Error; err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating client")
		return
	}

	respondJSON(w, http.StatusCreated, client)
}

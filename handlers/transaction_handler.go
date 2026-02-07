package handlers

// handlers menangani HTTP requests dan responses dari client

import (
	"category-api-ss2/models"
	"category-api-ss2/services"
	"encoding/json"
	"net/http"
)

// TransactionHandler menangani HTTP requests untuk transaksi
type TransactionHandler struct {
	service *services.TransactionService
}

// NewTransactionHandler membuat TransactionHandler dengan injected service
func NewTransactionHandler(service *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// HandleCheckout menangani POST ke /api/checkout
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Checkout memproses pembelian produk dengan update stok otomatis
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

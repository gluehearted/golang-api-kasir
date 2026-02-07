package handlers

import (
	"category-api-ss2/services"
	"encoding/json"
	"net/http"
	"time"
)

// ReportHandler menangani HTTP requests untuk laporan penjualan
type ReportHandler struct {
	service *services.ReportService
}

// NewReportHandler membuat ReportHandler dengan injected service
func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// HandleReport menangani GET ke /api/report
func (h *ReportHandler) HandleReport(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetReport(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetReport mengambil laporan penjualan hari ini atau dalam range tanggal
func (h *ReportHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Check jika "hari-ini"
	if path == "/api/report/hari-ini" {
		report, err := h.service.GetDailySalesReport()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
		return
	}

	// Check query parameter untuk date range
	startStr := r.URL.Query().Get("start_date")
	endStr := r.URL.Query().Get("end_date")

	if startStr != "" && endStr != "" {
		startDate, err := time.Parse("2006-01-02", startStr)
		if err != nil {
			http.Error(w, "Invalid start_date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endStr)
		if err != nil {
			http.Error(w, "Invalid end_date format (use YYYY-MM-DD)", http.StatusBadRequest)
			return
		}

		report, err := h.service.GetSalesReportByDateRange(startDate, endDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
		return
	}

	http.Error(w, "Use /api/report/hari-ini or /api/report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD", http.StatusBadRequest)
}

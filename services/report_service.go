package services

// services berisi business logic dan validasi untuk setiap domain

import (
	"category-api-ss2/models"
	"category-api-ss2/repositories"
	"time"
)

// ReportService menangani business logic untuk laporan penjualan
type ReportService struct {
	repo *repositories.ReportRepository
}

// NewReportService membuat ReportService dengan injected repository
func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

// GetDailySalesReport mengambil laporan penjualan hari ini
func (s *ReportService) GetDailySalesReport() (*models.SalesReport, error) {
	return s.repo.GetDailySalesReport()
}

// GetSalesReportByDateRange mengambil laporan penjualan dalam range tanggal
func (s *ReportService) GetSalesReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	return s.repo.GetSalesReportByDateRange(startDate, endDate)
}

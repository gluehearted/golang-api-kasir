package repositories

// repositories adalah data access layer untuk database operations

import (
	"category-api-ss2/models"
	"database/sql"
	"time"
)

// ReportRepository menangani operasi database untuk laporan penjualan
type ReportRepository struct {
	db *sql.DB
}

// NewReportRepository membuat ReportRepository dengan injected database connection
func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// GetDailySalesReport mengambil laporan penjualan hari ini
func (repo *ReportRepository) GetDailySalesReport() (*models.SalesReport, error) {
	today := time.Now().Format("2006-01-02")
	// tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	// Total revenue dan transaksi hari ini
	var totalRevenue, totalTransaksi int
	query := `
		SELECT COALESCE(SUM(total_amount), 0) as revenue, COUNT(*) as transaksi
		FROM transactions
		WHERE DATE(created_at) = $1
	`
	err := repo.db.QueryRow(query, today).Scan(&totalRevenue, &totalTransaksi)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Produk terlaris hari ini
	var produkTerlaris *models.TopProduct
	queryTop := `
		SELECT p.name, SUM(td.quantity) as qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) = $1
		GROUP BY p.id, p.name
		ORDER BY qty DESC
		LIMIT 1
	`
	var nama string
	var qty int
	err = repo.db.QueryRow(queryTop, today).Scan(&nama, &qty)
	if err == nil {
		produkTerlaris = &models.TopProduct{Nama: nama, QtyTerjual: qty}
	}

	return &models.SalesReport{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: produkTerlaris,
		StartDate:      time.Now(),
		EndDate:        time.Now(),
	}, nil
}

// GetSalesReportByDateRange mengambil laporan penjualan dalam range tanggal
func (repo *ReportRepository) GetSalesReportByDateRange(startDate, endDate time.Time) (*models.SalesReport, error) {
	// Total revenue dan transaksi dalam range
	var totalRevenue, totalTransaksi int
	query := `
		SELECT COALESCE(SUM(total_amount), 0) as revenue, COUNT(*) as transaksi
		FROM transactions
		WHERE created_at >= $1 AND created_at <= $2
	`
	err := repo.db.QueryRow(query, startDate, endDate.AddDate(0, 0, 1)).Scan(&totalRevenue, &totalTransaksi)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Produk terlaris dalam range
	var produkTerlaris *models.TopProduct
	queryTop := `
		SELECT p.name, SUM(td.quantity) as qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at >= $1 AND t.created_at <= $2
		GROUP BY p.id, p.name
		ORDER BY qty DESC
		LIMIT 1
	`
	var nama string
	var qty int
	err = repo.db.QueryRow(queryTop, startDate, endDate.AddDate(0, 0, 1)).Scan(&nama, &qty)
	if err == nil {
		produkTerlaris = &models.TopProduct{Nama: nama, QtyTerjual: qty}
	}

	return &models.SalesReport{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: produkTerlaris,
		StartDate:      startDate,
		EndDate:        endDate,
	}, nil
}

package models

// models berisi data structures untuk entities aplikasi

import "time"

// SalesReport mewakili ringkasan penjualan untuk periode tertentu
type SalesReport struct {
	TotalRevenue   int         `json:"total_revenue"`
	TotalTransaksi int         `json:"total_transaksi"`
	ProdukTerlaris *TopProduct `json:"produk_terlaris,omitempty"`
	StartDate      time.Time   `json:"start_date,omitempty"`
	EndDate        time.Time   `json:"end_date,omitempty"`
}

// TopProduct mewakili produk dengan penjualan tertinggi
type TopProduct struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

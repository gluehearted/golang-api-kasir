package models

// models berisi data structures untuk entities aplikasi

// Product mewakili entity produk
type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID *int   `json:"category_id,omitempty"`
}

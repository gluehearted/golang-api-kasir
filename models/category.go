package models

// models berisi data structures untuk entities aplikasi

// Category mewakili entity kategori
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

package models

// models berisi data structures untuk entities aplikasi

import "time"

// Transaction mewakili entity transaksi penjualan
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

// TransactionDetail mewakili detail item dalam transaksi
type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"`
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

// CheckoutItem adalah item dalam request checkout
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

// CheckoutRequest adalah wrapper untuk request body checkout
type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}

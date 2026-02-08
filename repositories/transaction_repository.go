package repositories

import (
	"category-api-ss2/models"
	"database/sql"
	"fmt"
)

// TransactionRepository menangani operasi database untuk transaksi
type TransactionRepository struct {
	db *sql.DB
}

// NewTransactionRepository membuat TransactionRepository dengan injected database connection
func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction membuat transaksi baru dengan atomic database operation
func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	// Begin atomic transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Validate products dan calculate totals
	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice int
		var productName string

		// Fetch product details
		err := tx.QueryRow("SELECT name, price FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		// Calculate subtotal dan update total
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// Update product stock
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Add to details
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// Insert transaction record
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// INSERT TRANSACTION DETAILS
	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		args := []interface{}{}

		for i, detail := range details {
			if i > 0 {
				query += ", "
			}
			query += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
			args = append(args, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
			detail.TransactionID = transactionID
			details[i] = detail
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Return transaction with all details
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

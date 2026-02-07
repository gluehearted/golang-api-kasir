package repositories

import (
	"category-api-ss2/models"
	"database/sql"
	"fmt"
)

// TransactionRepository struct menyimpan reference ke database connection
type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CreateTransaction membuat transaksi baru dengan atomic operation
// Validasi produk, hitung harga, update stok, dan insert record

//   - Fetch product details (name, price, stock)
//   - Validate: produk harus ada
//   - Calculate subtotal: price * quantity
//   - UPDATE product stock: stock - quantity
//   - Add ke details slice

// 3. INSERT transaction record (dengan total_amount)
// 4. LOOP untuk insert setiap transaction_detail
// 5. COMMIT transaction - jika semua sukses
// 6. ROLLBACK jika ada error di tengah jalan (via defer)

// ERROR HANDLING
// Jika ada error di tengah proses (contoh: produk tidak ada, stock tidak cukup):
// - Defer rollback akan otomatis merollback semua changes
// - Return error ke caller
// - Database state tidak berubah (atomic)
func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	// BEGIN transaction - mulai atomic operation
	// Ini memastikan semua database operations succeed atau semua rollback
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Defer rollback jika ada error

	// VALIDATE & CALCULATE
	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// Loop untuk setiap item yang dibeli
	for _, item := range items {
		var productPrice, stock int
		var productName string

		// Fetch product details dari database
		// Menggunakan tx (transaction) bukan repo.db untuk atomicity
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			// Produk tidak ditemukan - return error
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			// Other database errors
			return nil, err
		}

		// Calculate subtotal: harga x jumlah
		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// UPDATE product stock - kurangi stok karena terjual
		// Menggunakan tx untuk consistency dalam transaction
		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		// Append detail ke details slice
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// INSERT TRANSACTION RECORD
	// Insert transaksi utama dan dapatkan ID yang di-generate
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// INSERT TRANSACTION DETAILS
	// Loop untuk insert detail transaksi untuk setiap item
	for i := range details {
		// Set transaction_id yang baru di-generate
		details[i].TransactionID = transactionID

		// Insert transaction_detail record
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	// COMMIT
	// Jika semua berhasil, commit transaction
	// Ini akan membuat semua changes permanent
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// RETURN RESULT
	// Return Transaction object dengan semua details
	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

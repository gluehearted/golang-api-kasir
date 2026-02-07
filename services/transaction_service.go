package services

// services berisi business logic dan validasi untuk setiap domain

import (
	"category-api-ss2/models"
	"category-api-ss2/repositories"
)

// TransactionService menangani business logic untuk transaksi
type TransactionService struct {
	repo *repositories.TransactionRepository
}

// NewTransactionService membuat TransactionService dengan injected repository
func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

// Checkout memproses pembelian produk dengan update stok otomatis
func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items)
}

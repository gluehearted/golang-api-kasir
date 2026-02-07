package services

// services berisi business logic dan validasi untuk setiap domain

import (
	"category-api-ss2/models"
	"category-api-ss2/repositories"
)

// ProductService menangani business logic untuk produk
type ProductService struct {
	repo *repositories.ProductRepository
}

// NewProductService membuat ProductService dengan injected repository
func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// GetAll mengambil semua produk dengan optional filter by name
func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.repo.GetAll(name)
}

// Create membuat produk baru
func (s *ProductService) Create(data *models.Product) error {
	return s.repo.Create(data)
}

// GetByID mengambil detail produk berdasarkan ID
func (s *ProductService) GetByID(id int) (*models.Product, error) {
	return s.repo.GetByID(id)
}

// Update memperbarui data produk
func (s *ProductService) Update(product *models.Product) error {
	return s.repo.Update(product)
}

// Delete menghapus produk berdasarkan ID
// Parameter id adalah product ID yang ingin dihapus
func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

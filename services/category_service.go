package services

// services berisi business logic dan validasi untuk setiap domain

import (
	"category-api-ss2/models"
	"category-api-ss2/repositories"
)

// CategoryService menangani business logic untuk kategori
type CategoryService struct {
	repo *repositories.CategoryRepository
}

// NewCategoryService membuat CategoryService dengan injected repository
func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

// GetAll mengambil semua kategori
func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.repo.GetAll()
}

// Create membuat kategori baru
func (s *CategoryService) Create(data *models.Category) error {
	return s.repo.Create(data)
}

// GetByID mengambil detail kategori berdasarkan ID
func (s *CategoryService) GetByID(id int) (*models.Category, error) {
	return s.repo.GetByID(id)
}

// Update memperbarui data kategori
func (s *CategoryService) Update(category *models.Category) error {
	return s.repo.Update(category)
}

// Delete menghapus kategori berdasarkan ID
func (s *CategoryService) Delete(id int) error {
	return s.repo.Delete(id)
}

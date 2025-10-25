package services

import (
	"evernos-api2/models"
	"evernos-api2/repositories"
	"errors"
	"strings"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

// GetCategoryByID retrieves a category by ID
func (s *CategoryService) GetCategoryByID(id uint) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(namaCategory string) (*models.Category, error) {
	// Validate input
	if strings.TrimSpace(namaCategory) == "" {
		return nil, errors.New("nama category tidak boleh kosong")
	}

	category := &models.Category{
		NamaCategory: strings.TrimSpace(namaCategory),
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// UpdateCategory updates an existing category
func (s *CategoryService) UpdateCategory(id uint, namaCategory string) (*models.Category, error) {
	// Validate input
	if strings.TrimSpace(namaCategory) == "" {
		return nil, errors.New("nama category tidak boleh kosong")
	}

	// Check if category exists
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("category tidak ditemukan")
	}

	// Update category
	category.NamaCategory = strings.TrimSpace(namaCategory)
	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category by ID
func (s *CategoryService) DeleteCategory(id uint) error {
	// Check if category exists
	exists, err := s.categoryRepo.CheckExists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("category tidak ditemukan")
	}

	return s.categoryRepo.Delete(id)
}
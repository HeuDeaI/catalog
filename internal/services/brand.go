package services

import (
	"catalog/internal/models"
	"catalog/internal/repositories"
)

type BrandService interface {
	CreateBrand(brand *models.Brand) error
	GetBrandByID(id uint) (*models.Brand, error)
	GetAllBrands() ([]*models.Brand, error)
	UpdateBrand(brand *models.Brand) error
	DeleteBrand(id uint) error
}

type brandService struct {
	repo repositories.BrandRepository
}

func NewBrandService(repo repositories.BrandRepository) BrandService {
	return &brandService{repo: repo}
}

func (s *brandService) CreateBrand(brand *models.Brand) error {
	return s.repo.Create(brand)
}

func (s *brandService) GetBrandByID(id uint) (*models.Brand, error) {
	return s.repo.GetByID(id)
}

func (s *brandService) GetAllBrands() ([]*models.Brand, error) {
	return s.repo.GetAll()
}

func (s *brandService) UpdateBrand(brand *models.Brand) error {
	return s.repo.Update(brand)
}

func (s *brandService) DeleteBrand(id uint) error {
	return s.repo.Delete(id)
}

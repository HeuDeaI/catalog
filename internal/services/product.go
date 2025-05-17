package services

import (
	"catalog/internal/models"
	"catalog/internal/repositories"

	"gorm.io/gorm"
)

type ProductService interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uint) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
}

type productService struct {
	repo repositories.ProductRepository
	db   *gorm.DB
}

func NewProductService(repo repositories.ProductRepository, db *gorm.DB) ProductService {
	return &productService{repo: repo, db: db}
}

func (s *productService) CreateProduct(product *models.Product) error {
	if err := s.repo.Create(product); err != nil {
		return err
	}

	if len(product.SkinTypeIDs) > 0 {
		var skinTypes []models.SkinType
		if err := s.db.Where("id IN ?", product.SkinTypeIDs).Find(&skinTypes).Error; err != nil {
			return err
		}
		if err := s.db.Model(product).Association("SkinTypes").Replace(skinTypes); err != nil {
			return err
		}
	}

	return nil
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) GetAllProducts() ([]*models.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) UpdateProduct(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}

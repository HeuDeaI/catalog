package services

import (
	"catalog/internal/models"
	"catalog/internal/repositories"
	"fmt"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(product *models.Product, filePath string) error
	GetProductByID(id uint) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(product *models.Product, filePath string) error {
	objectName := fmt.Sprintf("%s.png", uuid.New().String())
	imageURL, err := s.repo.UploadImage(objectName, filePath)
	if err != nil {
		return err
	}
	product.ImageURL = imageURL
	return s.repo.Create(product)
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

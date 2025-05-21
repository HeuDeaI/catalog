package services

import (
	"catalog/internal/models"
	"catalog/internal/repositories"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(product *models.Product, filePath string) error
	GetProductByID(id uint) (*models.Product, error)
	GetAllProducts() ([]*models.Product, error)
	GetProductsByFilter(minPrice float64, maxPrice float64) ([]*models.Product, error)
	UpdateProduct(product *models.Product, filePath string) error
	DeleteProduct(id uint) error
}

type productService struct {
	repo      repositories.ProductRepository
	imageRepo repositories.ProductImageRepository
}

func NewProductService(repo repositories.ProductRepository, imageRepo repositories.ProductImageRepository) ProductService {
	return &productService{repo: repo, imageRepo: imageRepo}
}

func (s *productService) CreateProduct(product *models.Product, filePath string) error {
	objectName := fmt.Sprintf("%s.png", uuid.New().String())
	imageURL, err := s.imageRepo.UploadImage(filePath, objectName)
	if err != nil {
		return err
	}
	product.ImageURL = imageURL
	_ = os.Remove(filePath)
	return s.repo.Create(product)
}

func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *productService) GetAllProducts() ([]*models.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) UpdateProduct(product *models.Product, filePath string) error {
	if filePath != "" {
		if product.ImageURL != "" {
			parts := strings.Split(product.ImageURL, "/")
			oldObjectName := parts[len(parts)-1]
			if err := s.imageRepo.DeleteImage(oldObjectName); err != nil {
				return err
			}
		}

		objectName := fmt.Sprintf("%s.png", uuid.New().String())
		imageURL, err := s.imageRepo.UploadImage(filePath, objectName)
		if err != nil {
			return err
		}
		product.ImageURL = imageURL
		_ = os.Remove(filePath)
	}
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id uint) error {
	product, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if product.ImageURL != "" {
		parts := strings.Split(product.ImageURL, "/")
		objectName := parts[len(parts)-1]
		if err := s.imageRepo.DeleteImage(objectName); err != nil {
			return err
		}
	}
	return s.repo.Delete(id)
}

func (s *productService) GetProductsByFilter(minPrice float64, maxPrice float64) ([]*models.Product, error) {
	return s.repo.GetByFilter(minPrice, maxPrice)
}

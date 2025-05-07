package repositories

import "catalog/internal/models"

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetAll() ([]*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) Create(product *models.Product) error {
	return nil
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	return nil, nil
}

func (r *productRepository) GetAll() ([]*models.Product, error) {
	return nil, nil
}

func (r *productRepository) Update(product *models.Product) error {
	return nil
}

func (r *productRepository) Delete(id uint) error {
	return nil
}

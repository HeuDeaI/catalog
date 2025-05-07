package repositories

import "catalog/internal/models"

type BrandRepository interface {
	Create(brand *models.Brand) error
	GetByID(id uint) (*models.Brand, error)
	GetAll() ([]*models.Brand, error)
	Update(brand *models.Brand) error
	Delete(id uint) error
}

type brandRepository struct{}

func NewBrandRepository() BrandRepository {
	return &brandRepository{}
}

func (r *brandRepository) Create(brand *models.Brand) error {
	return nil
}

func (r *brandRepository) GetByID(id uint) (*models.Brand, error) {
	return nil, nil
}

func (r *brandRepository) GetAll() ([]*models.Brand, error) {
	return nil, nil
}

func (r *brandRepository) Update(brand *models.Brand) error {
	return nil
}

func (r *brandRepository) Delete(id uint) error {
	return nil
}

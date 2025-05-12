package repositories

import (
	"catalog/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetAll() ([]*models.Product, error)
	GetBySkinTypeID(skinTypeID uint) ([]*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Brand").Preload("SkinType").First(&product, id).Error
	return &product, err
}

func (r *productRepository) GetAll() ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Preload("Brand").Preload("SkinType").Find(&products).Error
	return products, err
}

func (r *productRepository) GetBySkinTypeID(skinTypeID uint) ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Preload("Brand").Preload("SkinType").Where("skin_type_id = ?", skinTypeID).Find(&products).Error
	return products, err
}

func (r *productRepository) Update(product *models.Product) error {
	err := r.db.Model(&product).Updates(product).Error
	if err == nil {
		r.db.Preload("Brand").Preload("SkinType").First(&product, product.ID)
	}
	return err
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

package repositories

import (
	"catalog/internal/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetAll() ([]*models.Product, error)
	GetByFilter(minPrice, maxPrice float64) ([]*models.Product, error)
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
	if err := r.db.Create(product).Error; err != nil {
		return err
	}

	if len(product.SkinTypeIDs) > 0 {
		var skinTypes []models.SkinType
		if err := r.db.Where("id IN ?", product.SkinTypeIDs).Find(&skinTypes).Error; err != nil {
			return err
		}
		if err := r.db.Model(product).Association("SkinTypes").Replace(skinTypes); err != nil {
			return err
		}
	}

	return nil
}

func (r *productRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Brand").Preload("SkinTypes").First(&product, id).Error
	return &product, err
}

func (r *productRepository) GetAll() ([]*models.Product, error) {
	var products []*models.Product
	err := r.db.Preload("Brand").Preload("SkinTypes").Find(&products).Error
	return products, err
}

func (r *productRepository) Update(product *models.Product) error {
	err := r.db.Model(&product).Updates(product).Error
	if err != nil {
		return err
	}

	if len(product.SkinTypeIDs) > 0 {
		var skinTypes []models.SkinType
		if err := r.db.Where("id IN ?", product.SkinTypeIDs).Find(&skinTypes).Error; err != nil {
			return err
		}
		if err := r.db.Model(product).Association("SkinTypes").Replace(skinTypes); err != nil {
			return err
		}
	}

	return r.db.Preload("Brand").Preload("SkinTypes").First(&product, product.ID).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) GetByFilter(minPrice, maxPrice float64) ([]*models.Product, error) {
	var products []*models.Product
	db := r.db.Preload("Brand").Preload("SkinTypes")

	if minPrice > 0 {
		db = db.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		db = db.Where("price <= ?", maxPrice)
	}

	err := db.Find(&products).Error
	return products, err
}

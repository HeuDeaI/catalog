package repositories

import (
	"catalog/internal/models"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	GetByID(id uint) (*models.Product, error)
	GetAll() ([]*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	UploadImage(objectName, filePath string) (string, error)
}

type productRepository struct {
	db       *gorm.DB
	minio    *minio.Client
	bucket   string
	minioURL string
}

func NewProductRepository(db *gorm.DB, minioClient *minio.Client, bucket string, minioURL string) ProductRepository {
	return &productRepository{db: db, minio: minioClient, bucket: bucket, minioURL: minioURL}
}

func (r *productRepository) UploadImage(objectName, filePath string) (string, error) {
	_, err := r.minio.FPutObject(context.Background(), r.bucket, objectName, filePath, minio.PutObjectOptions{ContentType: "image/png"})
	if err != nil {
		return "", err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", r.minioURL, r.bucket, objectName)
	return imageURL, nil
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
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
	if err == nil {
		r.db.Preload("Brand").Preload("SkinTypes").First(&product, product.ID)
	}
	return err
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

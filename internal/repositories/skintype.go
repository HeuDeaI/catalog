package repositories

import (
	"catalog/internal/models"

	"gorm.io/gorm"
)

type SkinTypeRepository interface {
	SeedSkinTypes()
}

type skinTypeRepository struct {
	db *gorm.DB
}

func NewSkinTypeRepository(db *gorm.DB) SkinTypeRepository {
	return &skinTypeRepository{db: db}
}

func (r *skinTypeRepository) SeedSkinTypes() {
	var count int64
	r.db.Model(&models.SkinType{}).Count(&count)
	if count == 0 {
		skinTypes := []models.SkinType{
			{ID: 1, Name: "Сухая"},
			{ID: 2, Name: "Жирная"},
			{ID: 3, Name: "Чувствительная"},
			{ID: 4, Name: "Комбинированная"},
			{ID: 5, Name: "Нормальная"},
		}
		r.db.Create(&skinTypes)
	}
}

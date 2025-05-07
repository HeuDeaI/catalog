package models

type Product struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	Name    string  `gorm:"type:varchar(255);not null" json:"name"`
	Price   float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	BrandID *uint   `gorm:"index" json:"brand_id"`
	Brand   Brand   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"brand"`
}

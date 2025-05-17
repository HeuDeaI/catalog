package models

type Product struct {
	ID          uint       `gorm:"primaryKey" json:"id" form:"id"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name" form:"name"`
	Price       float64    `gorm:"type:decimal(10,2);not null" json:"price" form:"price"`
	BrandID     *uint      `gorm:"index" json:"brand_id" form:"brand_id"`
	Brand       *Brand     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"brand"`
	SkinTypeIDs []uint     `gorm:"type:jsonb;serializer:json" json:"skin_type_ids" form:"skin_type_ids"`
	SkinTypes   []SkinType `gorm:"many2many:product_skin_types;" json:"skin_types"`
	ImageURL    string     `gorm:"type:varchar(1024);" json:"image_url"`
}

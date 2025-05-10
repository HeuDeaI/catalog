package models

type Brand struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(100);not null;unique" json:"name"`
	Products []Product `gorm:"foreignKey:BrandID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"products"`
}

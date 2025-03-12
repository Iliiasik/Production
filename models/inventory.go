package models

type Unit struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:varchar(50);not null"`
}

type RawMaterial struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"type:varchar(100);not null;unique"`
	UnitID      uint    `json:"unit_id" gorm:"index;not null"`
	Unit        Unit    `json:"unit" gorm:"foreignKey:UnitID;constraint:OnDelete:CASCADE"`
	Quantity    float64 `json:"quantity" gorm:"not null;default:0"`
	TotalAmount float64 `json:"total_amount" gorm:"not null;default:0"`
}

type FinishedGood struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"type:varchar(100);not null"`
	UnitID      uint    `json:"unit_id" gorm:"index;not null"`
	Unit        Unit    `json:"unit" gorm:"foreignKey:UnitID;constraint:OnDelete:CASCADE"`
	Quantity    float64 `json:"quantity" gorm:"not null;default:0"`
	TotalAmount float64 `json:"total_amount" gorm:"not null;default:0"`
}

type Ingredient struct {
	ID            uint         `json:"id" gorm:"primaryKey"`
	ProductID     uint         `json:"product_id" gorm:"index;not null"`
	Product       FinishedGood `json:"product" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	RawMaterialID uint         `json:"raw_material_id" gorm:"index;not null"`
	RawMaterial   RawMaterial  `json:"raw_material" gorm:"foreignKey:RawMaterialID;constraint:OnDelete:CASCADE"`
	Quantity      float64      `json:"quantity" gorm:"not null"`
}

func (Unit) TableName() string {
	return "units"
}

func (RawMaterial) TableName() string {
	return "raw_materials"
}

func (FinishedGood) TableName() string {
	return "finished_goods"
}

func (Ingredient) TableName() string {
	return "ingredients"
}

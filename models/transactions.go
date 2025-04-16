package models

import "time"

type RawMaterialPurchase struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	RawMaterialID uint        `json:"raw_material_id" gorm:"index;not null"`
	RawMaterial   RawMaterial `json:"raw_material" gorm:"foreignKey:RawMaterialID;constraint:OnDelete:CASCADE"`
	Quantity      float64     `json:"quantity" gorm:"not null"`
	TotalAmount   float64     `json:"total_amount" gorm:"not null"`
	PurchaseDate  time.Time   `json:"purchase_date" gorm:"not null"`
	EmployeeID    uint        `json:"employee_id" gorm:"index;not null"`
	Employee      Employee    `json:"employee" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

type ProductSale struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	ProductID   uint         `json:"product_id" gorm:"index;not null"`
	Product     FinishedGood `json:"product" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Quantity    float64      `json:"quantity" gorm:"not null"`
	TotalAmount float64      `json:"total_amount" gorm:"not null;default:0"`
	SaleDate    time.Time    `json:"sale_date" gorm:"not null"`
	EmployeeID  uint         `json:"employee_id" gorm:"index;not null"`
	Employee    Employee     `json:"employee" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

type ProductProduction struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	ProductID      uint         `json:"product_id" gorm:"index;not null"`
	Product        FinishedGood `json:"product" gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Quantity       float64      `json:"quantity" gorm:"not null"`
	ProductionDate time.Time    `json:"production_date" gorm:"not null"`
	EmployeeID     uint         `json:"employee_id" gorm:"index;not null"`
	Employee       Employee     `json:"employee" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
}

func (RawMaterialPurchase) TableName() string {
	return "raw_material_purchases"
}

func (ProductSale) TableName() string {
	return "product_sales"
}

func (ProductProduction) TableName() string {
	return "product_productions"
}

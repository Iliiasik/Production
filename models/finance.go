package models

type Budget struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	TotalAmount float64 `json:"total_amount" gorm:"not null"`
}


func (Budget) TableName() string {
	return "budgets"
}


package models

type Budget struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	TotalAmount float64 `json:"total_amount" gorm:"not null"`
	Markup      float64 `json:"markup" gorm:"not null;default:0"`
	SalaryBonus float64 `json:"salary_bonus" gorm:"not null;default:0"`
}

func (Budget) TableName() string {
	return "budgets"
}

package models

import "time"

type Budget struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	TotalAmount float64 `json:"total_amount" gorm:"not null"`
	Markup      float64 `json:"markup" gorm:"not null;default:0"`
	SalaryBonus float64 `json:"salary_bonus" gorm:"not null;default:0"`
}

func (Budget) TableName() string {
	return "budgets"
}

type Loan struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	BankName       string    `json:"bank_name" gorm:"type:varchar(100);not null"` // Название банка
	Amount         float64   `json:"amount" gorm:"not null"`                      // Сумма кредита
	InterestRate   float64   `json:"interest_rate" gorm:"not null"`               // Процент годовой
	MonthlyPayment float64   `json:"monthly_payment" gorm:"not null"`             // Ежемесячная выплата
	PenaltyTotal   float64   `json:"penalty_total" gorm:"not null;default:0"`     // Всего начислено пени
	IsRepaid       bool      `json:"is_repaid" gorm:"default:false"`              // Погашен ли кредит
	StartDate      time.Time `json:"start_date" gorm:"not null"`                  // Начало кредита
	Months         int       `json:"months" gorm:"not null"`                      // Срок в месяцах
}

type LoanPayment struct {
	ID       uint       `json:"id" gorm:"primaryKey"`
	LoanID   uint       `json:"loan_id" gorm:"index;not null"`
	DueDate  time.Time  `json:"due_date" gorm:"not null"`          // Плановая дата платежа
	Amount   float64    `json:"amount" gorm:"not null"`            // Плановая сумма (ежемесячная)
	Paid     bool       `json:"paid" gorm:"default:false"`         // Оплачено ли
	PaidDate *time.Time `json:"paid_date"`                         // Когда реально оплатили
	Penalty  float64    `json:"penalty" gorm:"not null;default:0"` // Пеня за просрочку (если есть)
}

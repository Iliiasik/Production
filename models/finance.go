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

type Credit struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Amount      float64         `json:"amount" gorm:"not null"`                  // Полученная сумма
	StartDate   time.Time       `json:"start_date" gorm:"not null"`              // Дата получения
	TermYears   int             `json:"term_years" gorm:"not null"`              // Срок в годах
	AnnualRate  float64         `json:"annual_rate" gorm:"not null"`             // % годовых
	PenaltyRate float64         `json:"penalty_rate" gorm:"not null"`            // Пеня (% в день)
	IsClosed    bool            `json:"is_closed" gorm:"not null;default:false"` // Погашен ли кредит
	Payments    []CreditPayment `json:"payments" gorm:"foreignKey:CreditID;constraint:OnDelete:CASCADE"`
}

func (Credit) TableName() string {
	return "credits"
}

type CreditPayment struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	CreditID         uint      `json:"credit_id" gorm:"index;not null"`
	MonthNumber      int       `json:"month_number" gorm:"not null"`       // № месяца
	PaymentDate      time.Time `json:"payment_date" gorm:"not null"`       // Дата выплаты
	PrincipalPart    float64   `json:"principal_part" gorm:"not null"`     // Часть кредита
	InterestPart     float64   `json:"interest_part" gorm:"not null"`      // Проценты
	TotalPayment     float64   `json:"total_payment" gorm:"not null"`      // Общая сумма платежа
	RemainingDebt    float64   `json:"remaining_debt" gorm:"not null"`     // Остаток кредита
	DaysOverdue      int       `json:"days_overdue" gorm:"not null"`       // Просрочка
	PenaltyAmount    float64   `json:"penalty_amount" gorm:"not null"`     // Пеня
	TotalWithPenalty float64   `json:"total_with_penalty" gorm:"not null"` // Итого к оплате
}

func (CreditPayment) TableName() string {
	return "credit_payments"
}

type PaymentsAggregate struct {
	SumPrincipalPart    float64 `json:"sum_principal_part"`
	SumInterestPart     float64 `json:"sum_interest_part"`
	SumTotalPayment     float64 `json:"sum_total_payment"`
	SumRemainingDebt    float64 `json:"sum_remaining_debt"`
	SumDaysOverdue      int64   `json:"sum_days_overdue"`
	SumPenaltyAmount    float64 `json:"sum_penalty_amount"`
	SumTotalWithPenalty float64 `json:"sum_total_with_penalty"`
}

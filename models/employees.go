package models

type Position struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:varchar(50);not null"`
}

type Employee struct {
	ID         uint     `json:"id" gorm:"primaryKey"`
	FullName   string   `json:"full_name" gorm:"type:varchar(100);not null"`
	PositionID uint     `json:"position_id" gorm:"index;not null"`
	Position   Position `json:"position" gorm:"foreignKey:PositionID;constraint:OnDelete:CASCADE"`
	Salary     float64  `json:"salary" gorm:"not null"`
	Address    string   `json:"address" gorm:"type:varchar(200)"`
	Phone      string   `json:"phone" gorm:"type:varchar(20)"`
}

type SalaryRecord struct {
	ID                 uint     `json:"id" gorm:"primaryKey"`
	EmployeeID         uint     `json:"employee_id" gorm:"index;not null"`
	Employee           Employee `json:"employee" gorm:"foreignKey:EmployeeID;constraint:OnDelete:CASCADE"`
	Year               int      `json:"year" gorm:"not null"`
	Month              int      `json:"month" gorm:"not null"`
	PurchaseCount      int      `json:"purchase_count" gorm:"not null;default:0"`
	ProductionCount    int      `json:"production_count" gorm:"not null;default:0"`
	SaleCount          int      `json:"sale_count" gorm:"not null;default:0"`
	TotalParticipation int      `json:"total_participation" gorm:"not null"`
	Bonus              float64  `json:"bonus" gorm:"not null"`
	TotalSalary        float64  `json:"total_salary" gorm:"not null"`
	IsPaid             bool     `json:"is_paid" gorm:"default:false"`
}

func (SalaryRecord) TableName() string {
	return "salary_records"
}

func (Position) TableName() string {
	return "positions"
}

func (Employee) TableName() string {
	return "employees"
}

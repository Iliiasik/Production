package models

import "time"

type Position struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"type:varchar(50);unique;not null"`
}

type Employee struct {
	ID                uint     `json:"id" gorm:"primaryKey"`
	FullName          string   `json:"full_name" gorm:"type:varchar(100);not null"`
	Username          string   `json:"username" gorm:"type:varchar(100);uniqueIndex"`
	PasswordHash      string   `json:"-"`
	IsPasswordChanged bool     `json:"is_password_changed" gorm:"default:false"`
	PositionID        uint     `json:"position_id" gorm:"index;not null"`
	Position          Position `json:"position" gorm:"foreignKey:PositionID;constraint:OnDelete:CASCADE"`
	Salary            float64  `json:"salary" gorm:"not null"`
	Address           string   `json:"address" gorm:"type:varchar(200)"`
	Phone             string   `json:"phone" gorm:"type:varchar(20)"`
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

type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(200);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"type:varchar(50);default:'Новая'"` // Новая, В работе, Завершена
	AssignedTo  uint      `json:"assigned_to" gorm:"not null"`                    // FK → employees(id)
	Employee    Employee  `json:"employee" gorm:"foreignKey:AssignedTo;references:ID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at"`
	DueDate     time.Time `json:"due_date"`
}

func (Task) TableName() string {
	return "tasks"
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

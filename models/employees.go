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

func (Position) TableName() string {
	return "positions"
}

func (Employee) TableName() string {
	return "employees"
}

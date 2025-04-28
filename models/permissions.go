package models

// Permission — атомарное право доступа
type Permission struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Name          string `json:"name" gorm:"unique;not null"`
	Description   string `json:"description" gorm:"type:text"`
	Category      string `json:"category" gorm:"type:varchar(50)"`
	VisibleToUser bool   `json:"visible_to_user"`
}

func (Permission) TableName() string {
	return "permissions"
}

// PositionPermission — связь между должностью и разрешением
type PositionPermission struct {
	PositionID   uint `json:"position_id" gorm:"primaryKey;not null"`
	PermissionID uint `json:"permission_id" gorm:"primaryKey;not null"`

	Position   Position   `json:"position" gorm:"foreignKey:PositionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (PositionPermission) TableName() string {
	return "position_permissions"
}

// UserPermission — индивидуальные права сотрудников с полным описанием прав
type UserPermission struct {
	EmployeeID   uint `json:"employee_id" gorm:"primaryKey;not null"`
	PermissionID uint `json:"permission_id" gorm:"primaryKey;not null"`

	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (UserPermission) TableName() string {
	return "user_permissions"
}

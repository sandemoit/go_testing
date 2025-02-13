package models

type UserRole struct {
	Id     uint `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	UserID uint `json:"user_id" gorm:"column:user_id"`
	RoleID uint `json:"role_id" gorm:"column:role_id"`
	Role   Role `gorm:"foreignKey:RoleID;references:id"`
	User   User `gorm:"foreignKey:UserID;references:id"`
}

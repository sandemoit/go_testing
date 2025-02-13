package models

type RolePermission struct {
	Id           uint       `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	PermissionId uint       `json:"permission_id" gorm:"column:permission_id"`
	Permission   Permission `gorm:"foreignKey:PermissionId"`
	RoleId       uint       `json:"role_id" gorm:"column:role_id"`
	Role         Role       `gorm:"foreignKey:RoleId"`
}

type RolePermissionRequest struct {
	PermissionId uint `json:"permission_id"`
	RoleId       uint `json:"role_id"`
}

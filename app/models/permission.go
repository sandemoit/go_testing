package models

type Permission struct {
	Id         uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Permission string `json:"permission" gorm:"unique;not null;column:permission;type:varchar(255)"`
	Roles      []Role `json:"roles" gorm:"many2many:role_permissions;"`
}

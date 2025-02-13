package models

import "time"

type User struct {
	Id            uint       `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name          string     `json:"name" form:"name" validate:"required,min=3" gorm:"type:varchar(255);not null"`
	Username      string     `json:"username" form:"username" validate:"required,min=3" gorm:"type:varchar(255);unique;not null"`
	Password      string     `json:"password,omitempty" form:"password" validate:"required,min=6" gorm:"type:varchar(255);"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	RememberToken string     `json:"remember_token" gorm:"type:varchar(255)"`
	Roles         []Role     `json:"roles,omitempty" gorm:"many2many:user_roles"`
}

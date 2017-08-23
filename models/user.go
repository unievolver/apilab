package models

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UserName string
	Email    string
	Password string
	Roles    []Role `gorm:"many2many:user_role"`
}

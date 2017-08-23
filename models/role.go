package models

import "github.com/jinzhu/gorm"

type Role struct {
	gorm.Model
	RoleName string
	Users    []User `gorm:"many2many:user_role"`
}

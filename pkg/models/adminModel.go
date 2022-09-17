package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Email    string
	Password string
}

type Tax struct {
	gorm.Model
	Category string
	Tax      int
}

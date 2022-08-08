package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	First_Name string
	Last_Name  string
	Email      string `gorm:"unique"`
	Password   string
	Status     bool
}

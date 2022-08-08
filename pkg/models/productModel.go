package models

import (
	"gorm.io/gorm"
)

type Products struct {
	gorm.Model
	Product_Name     string
	Brand_Name       string
	Product_Catogory string
	Description      string
	Color            string
	Size             int
	Stock            int
}

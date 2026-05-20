package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
}

package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title       string `gorm:"not null;type:varchar(255)"`
	Description string `gorm:"not null"`
}

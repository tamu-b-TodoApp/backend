package model

import "time"

type Todo struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string `gorm:"not null;type:varchar(255)"`
	Description string `gorm:"not null"`
}

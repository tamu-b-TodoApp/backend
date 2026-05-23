package model

import "time"

type Todo struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title" gorm:"not null;type:varchar(255)"`
	Description string    `json:"description" gorm:"not null"`
}

package model

import "time"

type User struct {
	ID              uint `gorm:"primarykey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Email           string `gorm:"not null;type:varchar(255);index"`
	Password        string `gorm:"not null;type:varchar(255)"`
	EmailVerifiedAt time.Time
}

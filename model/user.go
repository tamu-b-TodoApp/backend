package model

import "time"

type User struct {
	ID              uint       `json:"id" gorm:"primarykey"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Email           string     `json:"email" gorm:"not null;type:varchar(255);index"`
	Password        string     `json:"-" gorm:"not null;type:varchar(255)"`
	EmailVerifiedAt *time.Time `json:"email_verified_at" gorm:"default:null"`
}

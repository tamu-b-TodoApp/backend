package model

import "time"

type CompanyMember struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CompanyID uint      `json:"company_id" gorm:"not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
}

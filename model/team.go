package model

import "time"

type Team struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CompanyID uint      `json:"company_id" gorm:"not null;index"`
	Name      string    `json:"name" gorm:"not null;type:varchar(255)"`
}

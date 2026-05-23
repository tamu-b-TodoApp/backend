package model

import "time"

type RefreshToken struct {
	ID        uint      `gorm:"primarykey"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"not null;uniqueIndex;type:varchar(255)"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

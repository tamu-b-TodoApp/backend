package model

type TeamMember struct {
	TeamID          uint `gorm:"primaryKey"`
	CompanyMemberID uint `gorm:"primaryKey"`
}

package model

import "time"

type Todo struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	TeamID      uint       `json:"team_id" gorm:"not null;index"`
	ParentID    *uint      `json:"parent_id" gorm:"default:null"`
	AssigneeID  *uint      `json:"assignee_id" gorm:"default:null"`
	Title       string     `json:"title" gorm:"not null;type:varchar(255)"`
	Description string     `json:"description" gorm:"not null"`
	Status      string     `json:"status" gorm:"not null;type:varchar(20);default:'not_started'"`
	DueDate     *time.Time `json:"due_date" gorm:"type:date;default:null"`
	StoryPoints *int       `json:"story_points" gorm:"default:null"`
}

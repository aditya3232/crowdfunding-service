package entity

import "time"

type Campaign struct {
	ID               string          `gorm:"column:id;primaryKey"`
	UserID           string          `gorm:"column:user_id"` // 1 campaign belongs to 1 user
	Name             string          `gorm:"column:name"`
	ShortDescription string          `gorm:"column:short_description"`
	Description      string          `gorm:"column:description"`
	Perks            string          `gorm:"column:perks"`
	BackerCount      int             `gorm:"column:backer_count"`
	GoalAmount       int             `gorm:"column:goal_amount"`
	CurrentAmount    int             `gorm:"column:current_amount"`
	Slug             string          `gorm:"column:slug"`
	CreatedAt        time.Time       `gorm:"column:created_at"`
	UpdatedAt        time.Time       `gorm:"column:updated_at"`
	CampaignImages   []CampaignImage // 1 campaign has many campaign images
	User             User
}

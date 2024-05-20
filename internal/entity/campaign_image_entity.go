package entity

import "time"

type CampaignImage struct {
	ID         string    `gorm:"column:id;primaryKey"`
	CampaignID string    `gorm:"column:campaign_id"`
	FileName   string    `gorm:"column:file_name"`
	IsPrimary  int       `gorm:"column:is_primary"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (ci *CampaignImage) TableName() string {
	return "campaign_images"
}

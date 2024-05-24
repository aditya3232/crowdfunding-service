package entity

import "time"

type Transaction struct {
	ID         string    `gorm:"column:id;primaryKey"`
	CampaignID string    `gorm:"column:campaign_id"` // 1 transaction belongs to 1 campaign
	UserID     string    `gorm:"column:user_id"`     // 1 transaction belongs to 1 user
	Amount     int       `gorm:"column:amount"`
	Status     string    `gorm:"column:status"`
	Code       string    `gorm:"column:code"`
	PaymentURL string    `gorm:"column:payment_url"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	User       User
	Campaign   Campaign
}

func (t *Transaction) TableName() string {
	return "transactions"
}

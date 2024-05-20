package repository

import (
	"crowdfunding-service/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignImageRepository struct {
	Repository[entity.CampaignImage]
	Log *logrus.Logger
}

func NewCampaignImageRepository(log *logrus.Logger) *CampaignImageRepository {
	return &CampaignImageRepository{
		Log: log,
	}
}

// mark all campaign image as non primary
func (r *CampaignImageRepository) MarkAllAsNonPrimary(db *gorm.DB, campaignID string) error {
	return db.Model(new(entity.CampaignImage)).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error
}

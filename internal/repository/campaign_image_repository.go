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
// is_primary = false (artinya bukan 0, ya 1)
func (r *CampaignImageRepository) MarkAllAsNonPrimary(db *gorm.DB, campaignImage *entity.CampaignImage, campaignID string) (bool, error) {
	if err := db.Model(campaignImage).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error; err != nil {
		return false, err
	}

	return true, nil
}

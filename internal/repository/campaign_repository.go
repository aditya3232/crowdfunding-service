package repository

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignRepository struct {
	Repository[entity.Campaign]
	Log *logrus.Logger
}

func NewCampaignRepository(log *logrus.Logger) *CampaignRepository {
	return &CampaignRepository{
		Log: log,
	}
}

func (r *CampaignRepository) Search(db *gorm.DB, request *model.SearchCampaignRequest) ([]entity.Campaign, int64, error) {
	var campaigns []entity.Campaign

	// Apply filters and preload related data
	query := db.Scopes(r.FilterCampaign(request)).
		Preload("CampaignImages", "campaign_images.is_primary = 1").
		Preload("User")

	if err := query.
		Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&campaigns).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Campaign{}).Scopes(r.FilterCampaign(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
}

func (r *CampaignRepository) FilterCampaign(request *model.SearchCampaignRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if campaignName := request.CampaignName; campaignName != "" {
			campaignName = "%" + campaignName + "%"
			tx = tx.Where("campaigns.name LIKE ?", campaignName)
		}

		if userID := request.UserID; userID != "" {
			tx = tx.Where("campaigns.user_id = ?", userID)
		}

		// Filter based on user name
		if userName := request.UserName; userName != "" {
			userName = "%" + userName + "%"
			tx = tx.Joins("JOIN users ON users.id = campaigns.user_id").Where("users.name LIKE ?", userName)
		}

		return tx
	}
}

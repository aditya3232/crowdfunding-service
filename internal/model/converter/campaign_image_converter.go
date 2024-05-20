package converter

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
)

func CampaignImageToResponse(campaignImage *entity.CampaignImage) *model.CampaignImageResponse {
	return &model.CampaignImageResponse{
		ID:         campaignImage.ID,
		CampaignID: campaignImage.CampaignID,
		FileName:   campaignImage.FileName,
		IsPrimary:  campaignImage.IsPrimary,
	}
}

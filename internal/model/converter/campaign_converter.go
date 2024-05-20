package converter

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
)

func CampaignToResponse(campaign *entity.Campaign) *model.CampaignResponse {
	campaignImages := make([]model.CampaignImageResponse, len(campaign.CampaignImages))
	for i, campaignImage := range campaign.CampaignImages {
		campaignImages[i] = model.CampaignImageResponse(campaignImage)
	}

	user := model.UserResponse{
		ID:         campaign.User.ID,
		Name:       campaign.User.Name,
		Occupation: campaign.User.Occupation,
		Email:      campaign.User.Email,
		Avatar:     campaign.User.AvatarFileName,
		Role:       campaign.User.Role,
		CreatedAt:  campaign.User.CreatedAt,
		UpdatedAt:  campaign.User.UpdatedAt,
	}

	return &model.CampaignResponse{
		ID:               campaign.ID,
		UserID:           campaign.UserID,
		Name:             campaign.Name,
		ShortDescription: campaign.ShortDescription,
		Description:      campaign.Description,
		Perks:            campaign.Perks,
		BackerCount:      campaign.BackerCount,
		GoalAmount:       campaign.GoalAmount,
		CurrentAmount:    campaign.CurrentAmount,
		Slug:             campaign.Slug,
		CreatedAt:        campaign.CreatedAt,
		UpdatedAt:        campaign.UpdatedAt,
		CampaignImages:   campaignImages,
		User:             user,
	}
}

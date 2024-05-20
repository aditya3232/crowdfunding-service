package model

import "time"

type CampaignResponse struct {
	ID               string                  `json:"id"`
	UserID           string                  `json:"user_id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"short_description"`
	Description      string                  `json:"description"`
	Perks            string                  `json:"perks"`
	BackerCount      int                     `json:"backer_count"`
	GoalAmount       int                     `json:"goal_amount"`
	CurrentAmount    int                     `json:"current_amount"`
	Slug             string                  `json:"slug"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
	CampaignImages   []CampaignImageResponse `json:"campaign_images"`
	User             UserResponse            `json:"user"`
}

type CreateCampaignRequest struct {
	User             UserResponse // get from current user
	Name             string       `json:"name" validate:"required"`
	ShortDescription string       `json:"short_description" validate:"required"`
	Description      string       `json:"description" validate:"required"`
	Perks            string       `json:"perks" validate:"required"`
	GoalAmount       int          `json:"goal_amount" validate:"required"`
}

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
	CampaignImages   []CampaignImageResponse `json:"campaign_images,omitempty"`
	User             UserResponse            `json:"user,omitempty"`
}

type CreateCampaignRequest struct {
	UserID           string `json:"-" validate:"required,max=100,uuid"` // current user
	Name             string `json:"name" validate:"required,max=255"`
	ShortDescription string `json:"short_description" validate:"required,max=255"`
	Description      string `json:"description" validate:"required"`
	Perks            string `json:"perks" validate:"required"`
	GoalAmount       int    `json:"goal_amount" validate:"required,min=1"`
	Slug             string `json:"-"`
}

type UpdateCampaignRequest struct {
	ID               string `json:"-" validate:"required,max=100,uuid"`
	Name             string `json:"name" validate:"required,max=255"`
	ShortDescription string `json:"short_description" validate:"required,max=255"`
	Description      string `json:"description" validate:"required"`
	Perks            string `json:"perks" validate:"required"`
	GoalAmount       int    `json:"goal_amount" validate:"required,min=1"`
}

type SearchCampaignRequest struct {
	CampaignName string `json:"campaign_name" validate:"max=255"`
	UserID       string `json:"user_id" validate:"max=255"`
	UserName     string `json:"user_name" validate:"max=255"`
	Page         int    `json:"page" validate:"min=1"`
	Size         int    `json:"size" validate:"min=1,max=100"`
}

type GetCampaignRequest struct {
	ID string `json:"-" validate:"required,max=100,uuid"`
}

type DeleteCampaignRequest struct {
	ID string `json:"-" validate:"required,max=100,uuid"`
}

package converter

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
)

func CampaignTransactionToResponse(transaction *entity.Transaction) *model.CampaignTransactionResponse {
	return &model.CampaignTransactionResponse{
		ID:        transaction.ID,
		UserName:  transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
}

func UserTransactionToResponse(transaction *entity.Transaction) *model.UserTransactionResponse {
	campaign := model.Campaign{
		Name:     transaction.Campaign.Name,
		ImageURL: transaction.Campaign.CampaignImages[0].FileName,
	}

	return &model.UserTransactionResponse{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign:  campaign,
	}
}

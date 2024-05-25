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
	// pengecekan jika gambar tidak kosong, mengambil gambar pertama
	var imageURL string
	if len(transaction.Campaign.CampaignImages) > 0 {
		imageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	campaign := model.Campaign{
		Name:     transaction.Campaign.Name,
		ImageURL: imageURL,
	}

	return &model.UserTransactionResponse{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign:  campaign,
	}
}

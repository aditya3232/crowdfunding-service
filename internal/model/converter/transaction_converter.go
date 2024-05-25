package converter

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
)

func GetTransactionByCampaignIDToResponse(transaction *entity.Transaction) *model.GetTransactionByCampaignIDResponse {
	return &model.GetTransactionByCampaignIDResponse{
		ID:        transaction.ID,
		UserName:  transaction.User.Name,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}
}

func GetTransactionByUserIDToResponse(transaction *entity.Transaction) *model.GetTransactionByUserIDResponse {
	// pengecekan jika gambar tidak kosong, mengambil gambar pertama
	var imageURL string
	if len(transaction.Campaign.CampaignImages) > 0 {
		imageURL = transaction.Campaign.CampaignImages[0].FileName
	}

	campaign := model.Campaign{
		Name:     transaction.Campaign.Name,
		ImageURL: imageURL,
	}

	return &model.GetTransactionByUserIDResponse{
		ID:        transaction.ID,
		Amount:    transaction.Amount,
		Status:    transaction.Status,
		CreatedAt: transaction.CreatedAt,
		Campaign:  campaign,
	}
}

func TransactionToResponse(transaction *entity.Transaction) *model.TransactionResponse {
	return &model.TransactionResponse{
		ID:         transaction.ID,
		CampaignID: transaction.CampaignID,
		UserID:     transaction.UserID,
		Amount:     transaction.Amount,
		Status:     transaction.Status,
		Code:       transaction.Code,
		PaymentURL: transaction.PaymentURL,
		CreatedAt:  transaction.CreatedAt,
	}
}

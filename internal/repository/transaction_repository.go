package repository

import (
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	Repository[entity.Transaction]
	Log *logrus.Logger
}

func NewTransactionRepository(log *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{
		Log: log,
	}
}

func (r *TransactionRepository) GetTransactionByCampaignID(db *gorm.DB, request *model.GetTransactionByCampaignIDRequest) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction

	// Apply filters and preload related data
	if err := db.Scopes(r.FilterTransactionByCampaignID(request)).
		Preload("User").
		Order("created_at DESC").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Transaction{}).Scopes(r.FilterTransactionByCampaignID(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepository) FilterTransactionByCampaignID(request *model.GetTransactionByCampaignIDRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if campaignID := request.CampaignID; campaignID != "" {
			tx = tx.Where("campaign_id = ?", campaignID)
		}

		return tx
	}
}

func (r *TransactionRepository) GetTransactionByUserID(db *gorm.DB, request *model.GetTransactionByUserIDRequest) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction

	// Apply filters and preload related data
	if err := db.Scopes(r.FilterTransactionByUserID(request)).
		Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").
		Order("created_at DESC").
		Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Transaction{}).Scopes(r.FilterTransactionByUserID(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *TransactionRepository) FilterTransactionByUserID(request *model.GetTransactionByUserIDRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if userID := request.UserID; userID != "" {
			tx = tx.Where("user_id = ?", userID)
		}

		return tx
	}
}

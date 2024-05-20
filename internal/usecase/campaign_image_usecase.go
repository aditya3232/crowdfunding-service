package usecase

import (
	object_storing "crowdfunding-service/internal/gateway/object-storing"
	"crowdfunding-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignImageUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validate                *validator.Validate
	CampaignImageRepository *repository.CampaignImageRepository
	StoreObjectUseCase      *StoreObjectUseCase
	CampaignObject          *object_storing.StoreObject
}

func NewCampaignImageUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	campaignImageRepository *repository.CampaignImageRepository, storeObjectUseCase *StoreObjectUseCase, campaignObject *object_storing.StoreObject) *CampaignImageUseCase {
	return &CampaignImageUseCase{
		DB:                      db,
		Log:                     log,
		Validate:                validate,
		CampaignImageRepository: campaignImageRepository,
		StoreObjectUseCase:      storeObjectUseCase,
		CampaignObject:          campaignObject,
	}
}

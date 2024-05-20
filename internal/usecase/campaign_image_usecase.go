package usecase

import (
	"context"
	"crowdfunding-service/internal/entity"
	object_storing "crowdfunding-service/internal/gateway/object-storing"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/model/converter"
	"crowdfunding-service/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignImageUseCase struct {
	DB                      *gorm.DB
	Log                     *logrus.Logger
	Validate                *validator.Validate
	CampaignImageRepository *repository.CampaignImageRepository
	CampaignRepository      *repository.CampaignRepository
	UserRepository          *repository.UserRepository
	StoreObjectUseCase      *StoreObjectUseCase
	CampaignObject          *object_storing.StoreObject
}

func NewCampaignImageUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, campaignImageRepository *repository.CampaignImageRepository, campaignRepository *repository.CampaignRepository, userRepository *repository.UserRepository, storeObjectUseCase *StoreObjectUseCase, campaignObject *object_storing.StoreObject) *CampaignImageUseCase {
	return &CampaignImageUseCase{
		DB:                      db,
		Log:                     log,
		Validate:                validate,
		CampaignImageRepository: campaignImageRepository,
		CampaignRepository:      campaignRepository,
		UserRepository:          userRepository,
		StoreObjectUseCase:      storeObjectUseCase,
		CampaignObject:          campaignObject,
	}
}

func (u *CampaignImageUseCase) Create(ctx context.Context, request *model.CreateCampaignImageRequest) (*model.CampaignImageResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	// validasi image
	if request.FileImage == nil {
		if request.FileImage.Size > 5000000 {
			u.Log.Error("campaign image size is too large")
			return nil, fiber.ErrBadRequest
		}

		if !u.StoreObjectUseCase.IsImage(request.FileImage) {
			u.Log.Error("campaign image is not an image")
			return nil, fiber.ErrBadRequest
		}

		if !u.StoreObjectUseCase.IsValidImageFormat(request.FileImage) {
			u.Log.Error("campaign image is not a valid image format")
			return nil, fiber.ErrBadRequest
		}
	}

	// pengecekan campaign id
	campaign := new(entity.Campaign)
	if err := u.CampaignRepository.FindById(tx, campaign, request.CampaignID); err != nil {
		u.Log.WithError(err).Error("error finding campaign")
		return nil, fiber.ErrBadRequest
	}

	// pengecekan user id
	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.UserID); err != nil {
		u.Log.WithError(err).Error("error finding user")
		return nil, fiber.ErrBadRequest
	}

	// pengecekan apakah user yg login adalah pemilik campaign
	if campaign.UserID != request.UserID {
		u.Log.Error("user is not the owner of the campaign")
		return nil, fiber.ErrForbidden
	}

	// ubah image sebelumnya menjadi non primary (yang primary ditampilkan di halaman depan)
	// jika request is primary = true maka semua image di set menjadi non primary (kecuali yang baru di upload)
	// true -> 1, false -> 0
	// non primary = 0, primary = 1
	campaignImage := new(entity.CampaignImage)
	isPrimary := 0
	if request.IsPrimary {
		isPrimary = 1
		if _, err := u.CampaignImageRepository.MarkAllAsNonPrimary(tx, campaignImage, request.CampaignID); err != nil {
			u.Log.WithError(err).Error("error marking all campaign image as non primary")
			return nil, fiber.ErrInternalServerError
		}
	}

	// jika file image tidak kosong, set file name
	if request.FileImage != nil {
		campaignImage.FileName = "campaigns/image-" + uuid.New().String()
	}

	campaignImage.ID = uuid.New().String()
	campaignImage.CampaignID = request.CampaignID
	campaignImage.IsPrimary = isPrimary

	// setelah data campaign image di set, lakukan insert ke database
	if err := u.CampaignImageRepository.Create(tx, campaignImage); err != nil {
		u.Log.WithError(err).Error("error creating campaign image")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	// jika file image tidak kosong, simpan file image ke storage
	if request.FileImage != nil {
		if err := u.CampaignObject.StoreFromFileHeader(ctx, request.FileImage, campaignImage.FileName); err != nil {
			u.Log.WithError(err).Error("error storing object")
			return nil, fiber.ErrInternalServerError
		}
	}

	// menampilkan url untuk response
	if campaignImage.FileName != "" {
		campaignImage.FileName = u.CampaignObject.GetURLObject(campaignImage.FileName)
	}

	return converter.CampaignImageToResponse(campaignImage), nil

}

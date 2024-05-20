package usecase

import (
	"context"
	"crowdfunding-service/internal/entity"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/model/converter"
	"crowdfunding-service/internal/repository"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CampaignUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	CampaignRepository *repository.CampaignRepository
	UserRepository     *repository.UserRepository
}

func NewCampaignUseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate,
	campaignRepository *repository.CampaignRepository, userRepository *repository.UserRepository) *CampaignUseCase {
	return &CampaignUseCase{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		CampaignRepository: campaignRepository,
		UserRepository:     userRepository,
	}
}

func (u *CampaignUseCase) Create(ctx context.Context, request *model.CreateCampaignRequest) (*model.CampaignResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	// pengecekan user id dari current user
	user := new(entity.User)
	if err := u.UserRepository.FindById(tx, user, request.UserID); err != nil {
		u.Log.WithError(err).Error("error finding user by id")
		return nil, fiber.ErrBadRequest
	}

	// slug
	slugCandidate := fmt.Sprintf("%s-%s", request.Name, uuid.New().String())

	campaign := &entity.Campaign{
		ID:               uuid.New().String(),
		UserID:           request.UserID, // current user
		Name:             request.Name,
		ShortDescription: request.ShortDescription,
		Description:      request.Description,
		Perks:            request.Perks,
		GoalAmount:       request.GoalAmount,
		Slug:             slugCandidate, //jadi jgn tampilkan id campaign di url, tapi slug, agar seo friendly
	}

	if err := u.CampaignRepository.Create(tx, campaign); err != nil {
		u.Log.WithError(err).Error("error creating campaign")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	// add response user
	campaign.User = *user

	return converter.CampaignToResponse(campaign), nil

}

func (u *CampaignUseCase) Search(ctx context.Context, request *model.SearchCampaignRequest) ([]model.CampaignResponse, int64, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	campaigns, total, err := u.CampaignRepository.Search(tx, request)
	if err != nil {
		u.Log.WithError(err).Error("error searching campaign")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.CampaignResponse, len(campaigns))
	for i, campaign := range campaigns {
		responses[i] = *converter.CampaignToResponse(&campaign)
	}

	return responses, total, nil
}

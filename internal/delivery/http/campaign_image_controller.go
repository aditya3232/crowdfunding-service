package http

import (
	"crowdfunding-service/internal/delivery/http/middleware"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CampaignImageController struct {
	UseCase     *usecase.CampaignImageUseCase
	UserUseCase *usecase.UserUseCase
	Log         *logrus.Logger
}

func NewCampaignImageController(useCase *usecase.CampaignImageUseCase, userUseCase *usecase.UserUseCase, log *logrus.Logger) *CampaignImageController {
	return &CampaignImageController{
		UseCase:     useCase,
		UserUseCase: userUseCase,
		Log:         log,
	}
}

// create campaign image
func (c *CampaignImageController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateCampaignImageRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	auth := middleware.GetUser(ctx)
	responseGetUserByEmail, err := c.UserUseCase.GetByEmail(ctx.UserContext(), &model.GetUserByEmailRequest{Email: auth.Email})
	if err != nil {
		c.Log.WithError(err).Error("error getting user by email")
		return err
	}

	request.UserID = responseGetUserByEmail.ID

	file, err := ctx.FormFile("upload_campaign_image")
	if err != nil {
		c.Log.WithError(err).Error("error getting file")
		return fiber.ErrBadRequest
	}

	request.FileImage = file

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating campaign image")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CampaignImageResponse]{Data: response})
}

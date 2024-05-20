package http

import (
	"crowdfunding-service/internal/delivery/http/middleware"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CampaignController struct {
	UseCase     *usecase.CampaignUseCase
	UserUseCase *usecase.UserUseCase
	Log         *logrus.Logger
}

func NewCampaignController(useCase *usecase.CampaignUseCase, userUseCase *usecase.UserUseCase, log *logrus.Logger) *CampaignController {
	return &CampaignController{
		UseCase:     useCase,
		UserUseCase: userUseCase,
		Log:         log,
	}
}

func (c *CampaignController) CreateCampaign(ctx *fiber.Ctx) error {
	request := new(model.CreateCampaignRequest)
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

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating campaign")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.CampaignResponse]{Data: response})
}

func (c *CampaignController) List(ctx *fiber.Ctx) error {
	request := &model.SearchCampaignRequest{
		CampaignName: ctx.Query("campaign_name", ""),
		UserID:       ctx.Query("user_id", ""),
		UserName:     ctx.Query("user_name", ""),
		Page:         ctx.QueryInt("page", 1),
		Size:         ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching campaign")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.CampaignResponse]{
		Data:   responses,
		Paging: paging,
	})
}

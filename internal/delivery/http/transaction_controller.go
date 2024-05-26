package http

import (
	"crowdfunding-service/internal/delivery/http/middleware"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	UseCase     *usecase.TransactionUseCase
	UserUseCase *usecase.UserUseCase
	Log         *logrus.Logger
}

func NewTransactionController(useCase *usecase.TransactionUseCase, userUseCase *usecase.UserUseCase, log *logrus.Logger) *TransactionController {
	return &TransactionController{
		UseCase:     useCase,
		UserUseCase: userUseCase,
		Log:         log,
	}
}

// get transactions by campaign id
func (c *TransactionController) GetTransactionsByCampaignID(ctx *fiber.Ctx) error {
	request := &model.GetTransactionByCampaignIDRequest{
		CampaignID: ctx.Params("campaign_id"), // campaign_id diambil dari path parameter
		UserID:     ctx.Query("user_id", ""),
		Page:       ctx.QueryInt("page", 1),
		Size:       ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.GetTransactionsByCampaignID(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting transactions by campaign id")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.GetTransactionByCampaignIDResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// get transactions by user id
func (c *TransactionController) GetTransactionsByUserID(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	responseGetUserByEmail, err := c.UserUseCase.GetByEmail(ctx.UserContext(), &model.GetUserByEmailRequest{Email: auth.Email})
	if err != nil {
		c.Log.WithError(err).Error("error getting user by email")
		return err
	}

	request := &model.GetTransactionByUserIDRequest{
		UserID: responseGetUserByEmail.ID, // user_id diambil dari user yg login
		Page:   ctx.QueryInt("page", 1),
		Size:   ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.GetTransactionsByUserID(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting transactions by user id")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.GetTransactionByUserIDResponse]{
		Data:   responses,
		Paging: paging,
	})
}

// create transaction
func (c *TransactionController) CreateTransaction(ctx *fiber.Ctx) error {
	request := &model.CreateTransactionRequest{}
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

	response, err := c.UseCase.CreateTransaction(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating transaction")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.TransactionResponse]{Data: response})
}

// create transaction notification (get notification from midtrans)
func (c *TransactionController) CreateTransactionNotification(ctx *fiber.Ctx) error {
	request := &model.CreateTransactionNotificationRequest{}
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	err := c.UseCase.CreateTransactionNotification(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating transaction notification")
		return err
	}

	// return simple karena midtrans yg akses, cukup 200 OK
	return ctx.JSON(model.WebResponse[interface{}]{Data: "transaction notification created"})
}

package http

import (
	"crowdfunding-service/internal/delivery/http/middleware"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UseCase *usecase.UserUseCase
	Log     *logrus.Logger
}

func NewUserController(useCase *usecase.UserUseCase, log *logrus.Logger) *UserController {
	return &UserController{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *UserController) RegisterUser(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) List(ctx *fiber.Ctx) error {
	request := &model.SearchUserRequest{
		Name:  ctx.Query("name", ""),
		Email: ctx.Query("email", ""),
		Role:  ctx.Query("role", ""),
		Page:  ctx.QueryInt("page", 1),
		Size:  ctx.QueryInt("size", 10),
	}

	responses, total, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching user")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.UserResponse]{
		Data:   responses,
		Paging: paging,
	})
}

func (c *UserController) Get(ctx *fiber.Ctx) error {
	response, err := c.UseCase.Get(ctx.UserContext(), &model.GetUserRequest{ID: ctx.Params("userId")})
	if err != nil {
		c.Log.WithError(err).Error("error getting user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.ID = ctx.Params("userId")

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) CurrentUser(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	response, err := c.UseCase.GetByEmail(ctx.UserContext(), &model.GetUserByEmailRequest{Email: auth.Email})
	if err != nil {
		c.Log.WithError(err).Error("error getting current user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) UpdateAvatar(ctx *fiber.Ctx) error {
	request := new(model.UpdateAvatarRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	auth := middleware.GetUser(ctx)
	responseGetUserByEmail, err := c.UseCase.GetByEmail(ctx.UserContext(), &model.GetUserByEmailRequest{Email: auth.Email})
	if err != nil {
		c.Log.WithError(err).Error("error getting user by email")
		return err
	}

	request.ID = responseGetUserByEmail.ID

	file, err := ctx.FormFile("upload_avatar")
	if err != nil {
		c.Log.WithError(err).Error("error getting file")
		return fiber.ErrBadRequest
	}

	request.Avatar = file

	response, err := c.UseCase.UpdateAvatar(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating avatar")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

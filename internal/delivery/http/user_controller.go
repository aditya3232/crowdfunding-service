package http

import (
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
	request := &model.GetUserRequest{
		ID: ctx.Params("userId"),
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
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

func (c *UserController) Delete(ctx *fiber.Ctx) error {
	userId := ctx.Params("userId")

	request := &model.DeleteUserRequest{
		ID: userId,
	}

	err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error deleting user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

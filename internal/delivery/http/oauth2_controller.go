package http

import (
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Oauth2Controller struct {
	UseCase *usecase.Oauth2UseCase
	Log     *logrus.Logger
}

func NewOauth2Controller(useCase *usecase.Oauth2UseCase, log *logrus.Logger) *Oauth2Controller {
	return &Oauth2Controller{
		UseCase: useCase,
		Log:     log,
	}
}

func (c *Oauth2Controller) GoogleLogin(ctx *fiber.Ctx) error {
	response := c.UseCase.GoogleLogin()

	return ctx.JSON(model.WebResponse[*model.Oauth2GoogleLoginUrlResponse]{Data: response})
}

func (c *Oauth2Controller) GoogleCallback(ctx *fiber.Ctx) error {
	request := new(model.Oauth2GoogleCallbackRequest)
	if err := ctx.QueryParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request query")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.GoogleCallback(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error google callback")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.Oauth2GoogleCallbackResponse]{Data: response})
}

// revoke token
func (c *Oauth2Controller) RevokeToken(ctx *fiber.Ctx) error {
	request := new(model.Oauth2GoogleRevokeTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	err := c.UseCase.RevokeToken(request)
	if err != nil {
		c.Log.WithError(err).Error("error revoke token")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

// refresh token
func (c *Oauth2Controller) RefreshToken(ctx *fiber.Ctx) error {
	request := new(model.Oauth2GoogleRefreshTokenRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.RefreshToken(request)
	if err != nil {
		c.Log.WithError(err).Error("error refresh token")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.Oauth2GoogleRefreshTokenResponse]{Data: response})
}

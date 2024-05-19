package middleware

import (
	get_service "crowdfunding-service/internal/delivery/api-calling"
	"crowdfunding-service/internal/model"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type Oauth2Middleware struct {
	GoogleConfig           *oauth2.Config
	Config                 *viper.Viper
	Log                    *logrus.Logger
	GetOauth2GoogleService *get_service.GetOauth2GoogleService
}

func NewOauth2Middleware(googleConfig *oauth2.Config, config *viper.Viper, log *logrus.Logger, getOauth2GoogleService *get_service.GetOauth2GoogleService) *Oauth2Middleware {
	return &Oauth2Middleware{
		GoogleConfig:           googleConfig,
		Config:                 config,
		Log:                    log,
		GetOauth2GoogleService: getOauth2GoogleService,
	}
}

func (m *Oauth2Middleware) Oauth2Middleware(ctx *fiber.Ctx) error {
	request := &model.VerifyUserRequest{Token: ctx.Get("OAuth2-token", "NOT_FOUND")}
	// m.Log.Debugf("OAuth2-token : %s", request.Token)

	user, err := m.GetOauth2GoogleService.GetUserInfo(request.Token)
	if err != nil {
		m.Log.WithError(err).Error("failed to get user info")
		return fiber.ErrUnauthorized
	}

	userInfo := new(model.Oauth2GoogleCallbackResponse)
	if err := json.Unmarshal(user, userInfo); err != nil {
		m.Log.WithError(err).Error("error unmarshalling response")
		return fiber.ErrInternalServerError
	}

	userData := &model.Oauth2GoogleCallbackResponse{
		Email: userInfo.Email,
	}

	ctx.Locals("user", userData)

	return ctx.Next()
}

func GetUser(ctx *fiber.Ctx) *model.Oauth2GoogleCallbackResponse {
	return ctx.Locals("user").(*model.Oauth2GoogleCallbackResponse)
}

package api_calling

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type GetOauth2GoogleService struct {
	Log    *logrus.Logger
	Config *viper.Viper
}

func NewGetOauth2GoogleService(log *logrus.Logger, config *viper.Viper) *GetOauth2GoogleService {
	return &GetOauth2GoogleService{
		Log:    log,
		Config: config,
	}
}

// pemanggilan service google oauth2 & return data google oauth2
func (s *GetOauth2GoogleService) GetUserInfo(token string) ([]byte, error) {
	serviceGoogleOauth2 := s.Config.GetString("oauth2.google.getUserInfoWithAccessToken")
	client := fiber.AcquireClient()
	agent := client.Get(serviceGoogleOauth2 + token)

	defer fiber.ReleaseClient(client)

	status, response, errors := agent.Bytes()
	if status != fiber.StatusOK {
		if len(errors) > 0 {
			s.Log.WithError(errors[0]).Error("failed to get user info")
		} else {
			s.Log.Error("failed to get user info: unknown error")
		}
		return nil, fiber.ErrNotFound
	}

	return response, nil

}

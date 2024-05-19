package api_calling

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type SetOauth2GoogleService struct {
	Log    *logrus.Logger
	Config *viper.Viper
}

func NewSetOauth2GoogleService(log *logrus.Logger, config *viper.Viper) *SetOauth2GoogleService {
	return &SetOauth2GoogleService{
		Log:    log,
		Config: config,
	}
}

// revoke token
func (s *SetOauth2GoogleService) RevokeToken(token string) error {
	serviceGoogleOauth2 := s.Config.GetString("oauth2.google.revokeAccessToken")
	client := fiber.AcquireClient()
	agent := client.Post(serviceGoogleOauth2 + token)

	defer fiber.ReleaseClient(client)

	status, _, errors := agent.Bytes()
	if status != fiber.StatusOK {
		if len(errors) > 0 {
			s.Log.WithError(errors[0]).Error("failed to revoke token")
		} else {
			s.Log.Error("failed to revoke token: unknown error")
		}
		return fiber.ErrNotFound
	}

	return nil
}

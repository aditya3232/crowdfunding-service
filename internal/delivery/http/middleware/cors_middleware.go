package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CorsMiddleware struct {
	Config *viper.Viper
	Log    *logrus.Logger
}

func NewCorsMiddleware(config *viper.Viper, log *logrus.Logger) *CorsMiddleware {
	return &CorsMiddleware{
		Config: config,
		Log:    log,
	}
}

func (m *CorsMiddleware) CorsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: m.Config.GetString("cors.allowedOrigins"),
		AllowMethods: m.Config.GetString("cors.allowedMethods"),
		AllowHeaders: m.Config.GetString("cors.allowedHeaders"),
	})
}

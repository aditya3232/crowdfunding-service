package test

import (
	"crowdfunding-service/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

var (
	db           *gorm.DB
	app          *fiber.App
	log          *logrus.Logger
	validate     *validator.Validate
	viperConfig  *viper.Viper
	oauth2Google *oauth2.Config
	objectStore  *minio.Client
)

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	db = config.NewDatabase(viperConfig, log)
	oauth2Google = config.NewGoogleConfig(viperConfig)
	objectStore = config.NewMinio(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validate:    validate,
		Config:      viperConfig,
		Oauth2:      oauth2Google,
		ObjectStore: objectStore,
	})
}

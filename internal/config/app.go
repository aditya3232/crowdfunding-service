package config

import (
	"crowdfunding-service/internal/delivery/http"
	"crowdfunding-service/internal/delivery/http/route"
	"crowdfunding-service/internal/repository"
	"crowdfunding-service/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Log         *logrus.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	ObjectStore *minio.Client
}

func Bootstrap(config *BootstrapConfig) {
	// repositories
	userRepository := repository.NewUserRepository(config.Log)

	// usecases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)

	// controllers
	userController := http.NewUserController(userUseCase, config.Log)

	// routes
	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
	}
	routeConfig.Setup()
}

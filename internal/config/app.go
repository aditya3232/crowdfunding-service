package config

import (
	get_service "crowdfunding-service/internal/delivery/api-calling"
	"crowdfunding-service/internal/delivery/http"
	"crowdfunding-service/internal/delivery/http/middleware"
	"crowdfunding-service/internal/delivery/http/route"
	set_service "crowdfunding-service/internal/gateway/api-calling"
	object_storing "crowdfunding-service/internal/gateway/object-storing"
	payment_gateway "crowdfunding-service/internal/gateway/payment-gateway"
	"crowdfunding-service/internal/repository"
	"crowdfunding-service/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Log         *logrus.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	Oauth2      *oauth2.Config
	ObjectStore *minio.Client
}

func Bootstrap(config *BootstrapConfig) {
	// repositories
	userRepository := repository.NewUserRepository(config.Log)
	campaignRepository := repository.NewCampaignRepository(config.Log)
	campaignImageRepository := repository.NewCampaignImageRepository(config.Log)
	transactionRepository := repository.NewTransactionRepository(config.Log)

	// services
	GetOauth2GoogleService := get_service.NewGetOauth2GoogleService(config.Log, config.Config)
	SetOauth2GoogleService := set_service.NewSetOauth2GoogleService(config.Log, config.Config)

	// store objects
	objectStore := object_storing.NewUserObject(config.ObjectStore, config.Config, config.Log)

	// payment gateway
	midtransPaymentGateway := payment_gateway.NewMidtransPaymentGateway(config.Log, config.Config)

	// usecases
	oauth2UseCase := usecase.NewOauth2UseCase(config.DB, config.Log, config.Validate, config.Config, userRepository, config.Oauth2, GetOauth2GoogleService, SetOauth2GoogleService)
	storeObjectUseCase := usecase.NewStoreObjectUseCase(config.Log)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, storeObjectUseCase, objectStore)
	campaignUseCase := usecase.NewCampaignUseCase(config.DB, config.Log, config.Validate, campaignRepository, userRepository)
	campaignImageUseCase := usecase.NewCampaignImageUseCase(config.DB, config.Log, config.Validate, campaignImageRepository, campaignRepository, userRepository, storeObjectUseCase, objectStore)
	transactionUseCase := usecase.NewTransactionUseCase(config.DB, config.Log, config.Validate, transactionRepository, campaignRepository, userRepository, midtransPaymentGateway)

	// controllers
	oauth2Controller := http.NewOauth2Controller(oauth2UseCase, config.Log)
	userController := http.NewUserController(userUseCase, config.Log)
	campaignController := http.NewCampaignController(campaignUseCase, userUseCase, config.Log)
	campaignImageController := http.NewCampaignImageController(campaignImageUseCase, userUseCase, config.Log)
	transactionController := http.NewTransactionController(transactionUseCase, userUseCase, config.Log)

	// middleware
	oauth2Middleware := middleware.NewOauth2Middleware(config.Oauth2, config.Config, config.Log, GetOauth2GoogleService)
	corsMiddleware := middleware.NewCorsMiddleware(config.Config, config.Log)

	// routes
	routeConfig := route.RouteConfig{
		App:                     config.App,
		Oauth2Middleware:        oauth2Middleware,
		CorsMiddleware:          corsMiddleware,
		Oauth2Controller:        oauth2Controller,
		UserController:          userController,
		CampaignController:      campaignController,
		CampaignImageController: campaignImageController,
		TransactionController:   transactionController,
	}
	routeConfig.Setup()
}

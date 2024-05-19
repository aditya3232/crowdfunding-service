package usecase

import (
	"context"
	get_service "crowdfunding-service/internal/delivery/api-calling"
	"crowdfunding-service/internal/entity"
	set_service "crowdfunding-service/internal/gateway/api-calling"
	"crowdfunding-service/internal/model"
	"crowdfunding-service/internal/model/converter"
	"crowdfunding-service/internal/repository"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type Oauth2UseCase struct {
	DB                     *gorm.DB
	Log                    *logrus.Logger
	Validate               *validator.Validate
	Config                 *viper.Viper
	UserRepository         *repository.UserRepository
	GoogleConfig           *oauth2.Config
	GetOauth2GoogleService *get_service.GetOauth2GoogleService
	SetOauth2GoogleService *set_service.SetOauth2GoogleService
}

func NewOauth2UseCase(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, config *viper.Viper,
	userRepository *repository.UserRepository, googleConfig *oauth2.Config, getOauth2GoogleService *get_service.GetOauth2GoogleService,
	setOauth2GoogleService *set_service.SetOauth2GoogleService) *Oauth2UseCase {
	return &Oauth2UseCase{
		DB:                     db,
		Log:                    log,
		Validate:               validate,
		Config:                 config,
		UserRepository:         userRepository,
		GoogleConfig:           googleConfig,
		GetOauth2GoogleService: getOauth2GoogleService,
		SetOauth2GoogleService: setOauth2GoogleService,
	}
}

func (u *Oauth2UseCase) GoogleLogin() *model.Oauth2GoogleLoginUrlResponse {
	// gunakan oauth2.AccessTypeOffline agar bisa mendapatkan refresh token
	// refresh token akan expired setelah 6 bulan
	// refresh token akan expired jika user ganti password atau revoke akses
	// refresh token akan expired jika user hapus aplikasi dari google account
	loginUrl := &model.Oauth2GoogleLoginUrlResponse{
		LoginUrl: u.GoogleConfig.AuthCodeURL(u.Config.GetString("oauth2.google.stateString"), oauth2.AccessTypeOffline),
	}

	return converter.Oauth2GoogleLoginUrlToResponse(loginUrl)
}

func (u *Oauth2UseCase) GoogleCallback(ctx context.Context, request *model.Oauth2GoogleCallbackRequest) (*model.Oauth2GoogleCallbackResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	// check state
	if request.State != u.Config.GetString("oauth2.google.stateString") {
		u.Log.Error("invalid state")
		return nil, fiber.ErrUnauthorized
	}

	// Exchange authorization code for token
	token, err := u.GoogleConfig.Exchange(ctx, request.Code)
	if err != nil {
		u.Log.WithError(err).Error("failed to exchange token")
		return nil, fiber.ErrUnauthorized
	}

	// Retrieve user info using the token
	resp, err := u.GetOauth2GoogleService.GetUserInfo(token.AccessToken)
	if err != nil {
		u.Log.WithError(err).Error("failed to get user info")
		return nil, fiber.ErrNotFound
	}

	// parse response
	userInfo := new(model.Oauth2GoogleCallbackResponse)
	if err := json.Unmarshal(resp, userInfo); err != nil {
		u.Log.WithError(err).Error("error unmarshalling response")
		return nil, fiber.ErrInternalServerError
	}

	userData := &model.Oauth2GoogleCallbackResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Email:        userInfo.Email,
		ExpiresIn:    token.Expiry.String(),
	}

	// check user email
	user := &entity.User{Email: userInfo.Email}
	totalEmail, err := u.UserRepository.CountByEmail(tx, user)
	if err != nil {
		u.Log.WithError(err).Error("error checking email")
		return nil, fiber.ErrInternalServerError
	}

	if totalEmail == 0 {
		u.Log.Error("user not found")
		return nil, fiber.ErrNotFound
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.Oauth2GoogleCallbackToResponse(userData), nil
}

func (u *Oauth2UseCase) RevokeToken(request *model.Oauth2GoogleRevokeTokenRequest) error {
	accessToken := request.AccessToken

	if err := u.SetOauth2GoogleService.RevokeToken(accessToken); err != nil {
		u.Log.WithError(err).Error("failed to revoke access token")
		return fiber.ErrNotFound
	}

	return nil
}

func (u *Oauth2UseCase) RefreshToken(request *model.Oauth2GoogleRefreshTokenRequest) (*model.Oauth2GoogleRefreshTokenResponse, error) {
	if err := u.Validate.Struct(request); err != nil {
		u.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	token := &oauth2.Token{
		RefreshToken: request.RefreshToken,
	}

	newToken, err := u.GoogleConfig.TokenSource(context.Background(), token).Token()
	if err != nil {
		u.Log.WithError(err).Error("failed to refresh token")
		return nil, fiber.ErrUnauthorized
	}

	newAccessToken := &model.Oauth2GoogleRefreshTokenResponse{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresIn:    newToken.Expiry.String(),
	}

	return converter.Oauth2GoogleRefreshTokenToResponse(newAccessToken), nil

}

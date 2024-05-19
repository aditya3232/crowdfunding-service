package converter

import (
	"crowdfunding-service/internal/model"
)

func Oauth2GoogleCallbackToResponse(callback *model.Oauth2GoogleCallbackResponse) *model.Oauth2GoogleCallbackResponse {
	return &model.Oauth2GoogleCallbackResponse{
		AccessToken:  callback.AccessToken,
		RefreshToken: callback.RefreshToken,
		Email:        callback.Email,
		ExpiresIn:    callback.ExpiresIn,
	}
}

func Oauth2GoogleLoginUrlToResponse(url *model.Oauth2GoogleLoginUrlResponse) *model.Oauth2GoogleLoginUrlResponse {
	return &model.Oauth2GoogleLoginUrlResponse{
		LoginUrl: url.LoginUrl,
	}
}

func Oauth2GoogleRefreshTokenToResponse(refreshToken *model.Oauth2GoogleRefreshTokenResponse) *model.Oauth2GoogleRefreshTokenResponse {
	return &model.Oauth2GoogleRefreshTokenResponse{
		AccessToken:  refreshToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
		ExpiresIn:    refreshToken.ExpiresIn,
	}
}

package model

type Oauth2GoogleCallbackResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Email        string `json:"email"`
	ExpiresIn    string `json:"expires_in"`
}

type Oauth2GoogleLoginUrlResponse struct {
	LoginUrl string `json:"login_url"`
}

type Oauth2GoogleRefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

type Oauth2GoogleCallbackRequest struct {
	State string `json:"state" validate:"required"`
	Code  string `json:"code" validate:"required"` //pakai QueryParser
}

type Oauth2GoogleRevokeTokenRequest struct {
	AccessToken string `json:"access_token" validate:"required"`
}

type Oauth2GoogleRefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type VerifyUserRequest struct {
	Token string `validate:"required,max=100"`
}

type GetUserAuth struct {
	UserId string `json:"user_id"`
}

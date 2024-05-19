package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func NewGoogleConfig(config *viper.Viper) *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.GetString("oauth2.google.clientId"),
		ClientSecret: config.GetString("oauth2.google.clientSecret"),
		RedirectURL:  config.GetString("oauth2.google.redirectUrl"),
		Scopes: []string{config.GetString("oauth2.google.scopeUserEmail"),
			config.GetString("oauth2.google.scopeUserProfile")},
		Endpoint: google.Endpoint,
	}

	return conf

}

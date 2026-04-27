package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func NewGithubOauthConfig(redirectURL string, clientID string, clientSecret string) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}

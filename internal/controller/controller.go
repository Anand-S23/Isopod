package controller

import (
	"context"
	"encoding/json"
	"net/http"
)

type Controller struct {
    Ctx          context.Context
    Production   bool
    GithubOauth  oauth2.Config
    BaseURL      string
    CSRFToken    string
}

func NewController(
    ctx context.Context, 
    production bool, 
    baseURL string, 
    githubClientID string, 
    githubClientSecret string,
    csrfToken string,
) *Controller {
    githubOauthConfig := config.NewGithubOauthConfig(
        fmt.Sprintf("%s/auth/callback", baseURL),
        githubClientID,
        githubClientSecret,
    )

    return &Controller {
        Ctx: ctx,
        Production: production,
        GithubOauth: githubOauthConfig,
        BaseURL: baseURL,
        CSRFToken: csrfToken,
    }
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ErrMsg(message string) map[string]string {
    return map[string]string {"error": message}
}


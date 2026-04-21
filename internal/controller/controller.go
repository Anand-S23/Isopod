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

    store        *store.Store
}

func NewController(
    store *store.Store,
    env *config.Env,
    ctx context.Context, 
) *Controller {
    githubOauthConfig := config.NewGithubOauthConfig(
        fmt.Sprintf("%s/auth/callback", env.BaseURL),
        env.GithubClientID,
        env.GithubClientSecret,
    )

    return &Controller {
        Ctx: ctx,
        Production: env.Production,
        GithubOauth: githubOauthConfig,
        BaseURL: env.BaseURL,
        CSRFToken: env.CSRFToken,
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

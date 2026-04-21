package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anand-S23/isopod/internal/config"
	"github.com/Anand-S23/isopod/internal/store"
	"golang.org/x/oauth2"
)

type Controller struct {
	Ctx           context.Context
	Production    bool
	GithubOauth   oauth2.Config
	BaseURL       string
	CSRFToken     string
	EncryptionKey []byte

	store *store.Store
}

func NewController(
	store *store.Store,
	env *config.EnvVars,
	ctx context.Context,
) *Controller {
	githubOauthConfig := config.NewGithubOauthConfig(
		fmt.Sprintf("%s/auth/callback", env.BASE_URL),
		env.GITHUB_CLIENT_ID,
		env.GITHUB_CLIENT_SECRET,
	)

	return &Controller{
		Ctx:           ctx,
		Production:    env.PRODUCTION,
		GithubOauth:   githubOauthConfig,
		BaseURL:       env.BASE_URL,
		CSRFToken:     env.CSRF_TOKEN,
		EncryptionKey: env.ENCRYPTION_KEY,
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ErrMsg(message string) map[string]string {
	return map[string]string{"error": message}
}

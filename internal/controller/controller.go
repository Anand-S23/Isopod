package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anand-S23/isopod/internal/config"
	"github.com/Anand-S23/isopod/internal/store"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

type Controller struct {
	Ctx           context.Context
	Production    bool
	GithubOauth   oauth2.Config
	BaseURL       string
	AllowedOrigin string
	CSRFToken     string
	EncryptionKey []byte
	SessionStore  *sessions.CookieStore

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

	cookieStore := sessions.NewCookieStore(env.SESSION_KEY, env.SESSION_KEY)
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   env.PRODUCTION,
		SameSite: http.SameSiteLaxMode,
	}

	return &Controller{
		Ctx:           ctx,
		Production:    env.PRODUCTION,
		GithubOauth:   githubOauthConfig,
		BaseURL:       env.BASE_URL,
		AllowedOrigin: env.ALLOWED_ORIGIN,
		CSRFToken:     env.CSRF_TOKEN,
		EncryptionKey: env.ENCRYPTION_KEY,
		SessionStore:  cookieStore,
		store:         store,
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

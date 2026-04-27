package controller

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Anand-S23/isopod/internal/models"
)

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) error {
	url := c.GithubOauth.AuthCodeURL(c.CSRFToken)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func (c *Controller) Callback(w http.ResponseWriter, r *http.Request) error {
	state := r.FormValue("state")
	if state != c.CSRFToken {
		return WriteJSON(w, http.StatusBadRequest, ErrMsg("invalid CSRF token"))
	}

	code := r.FormValue("code")
	token, err := c.GithubOauth.Exchange(c.Ctx, code)
	if err != nil {
		errMsg := fmt.Sprintf("oauth exchange: %s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, ErrMsg(errMsg))
	}

	client := c.GithubOauth.Client(c.Ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		errMsg := fmt.Sprintf("get user: %s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, ErrMsg(errMsg))
	}
	defer resp.Body.Close()

	var ghUserData models.GHUserDto
	err = json.NewDecoder(resp.Body).Decode(&ghUserData)
	if err != nil {
		errMsg := fmt.Sprintf("decode user: %s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, ErrMsg(errMsg))
	}

	user, err := models.NewUser(ghUserData, c.EncryptionKey)
	if err != nil {
		errMsg := fmt.Sprintf("create user: %s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, ErrMsg(errMsg))
	}

	err = c.store.User.Upsert(c.Ctx, &user)
	if err != nil {
		errMsg := fmt.Sprintf("update or create user: %s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, ErrMsg(errMsg))
	}

	sess, err := c.getSession(r)
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}
	sess.Values[sessionValueUserID] = user.ID
	if err := sess.Save(r, w); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	http.Redirect(w, r, c.AllowedOrigin, http.StatusFound)
	return nil
}

func (c *Controller) Me(w http.ResponseWriter, r *http.Request) error {
	id, ok := c.currentUserID(r)
	if !ok {
		return WriteJSON(w, http.StatusUnauthorized, ErrMsg("not authenticated"))
	}
	u, err := c.store.User.GetByID(c.Ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return WriteJSON(w, http.StatusUnauthorized, ErrMsg("not authenticated"))
		}
		return WriteJSON(w, http.StatusInternalServerError, ErrMsg("load user"))
	}
	return WriteJSON(w, http.StatusOK, map[string]any{
		"id":            u.ID,
		"github_id":     u.GitHubID,
		"username":      u.Username,
		"email":         u.Email,
		"avatar_url":    u.AvatarURL,
		"created_at":    u.CreatedAt,
		"last_login_at": u.LastLoginAt,
	})
}

func (c *Controller) Logout(w http.ResponseWriter, r *http.Request) error {
	sess, err := c.getSession(r)
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}
	sess.Options.MaxAge = -1
	sess.Values = make(map[interface{}]interface{})
	if err := sess.Save(r, w); err != nil {
		return fmt.Errorf("save session: %w", err)
	}
	return WriteJSON(w, http.StatusOK, map[string]any{"ok": true})
}

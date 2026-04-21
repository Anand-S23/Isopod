package controller

import (
	"encoding/json"
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

	err = c.store.User.Upsert(&user)
	if err != nil {
		errMsg := fmt.Sprintf("update or create user: %s", err.Error())
		return WriteJSON(w, http.StatusBadRequest, ErrMsg(errMsg))
	}

	return nil
}

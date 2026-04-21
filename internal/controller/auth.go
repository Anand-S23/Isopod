package controller

func (c *Controller) Login(w http.ResponseWriter, r *http.Request) error {
	url := c.GithubOauth.AuthCodeURL(c.CSRFToken)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func (c *Controller) Callback(w http.ResponseWriter, r *http.Request) error {
	state := r.FormValue("state")
	if state != c.CSRFToken {
		return ErrMsg("invalid CSRF token")
	}

	code := r.FormValue("code")
	token, err := c.GithubOauth.Exchange(c.Ctx, code)
	if err != nil {
		return ErrMsg("oauth exchange: %w", err)
	}

	client := c.GithubOauth.Client(c.Ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return ErrMsg("get user: %w", err)
	}
	defer resp.Body.Close()

	var ghUserData models.GHUserDto
	err := json.NewDecoder(resp.Body).Decode(&ghUserData)
	if err != nil {
		return ErrMsg("decode user: %w", err)
	}

	user := models.NewUser(ghUserData, c.EncryptionKey)
	// TODO: Upsert user in database and create session
	return nil
}

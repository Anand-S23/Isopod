package controller

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const isopodSession = "isopod"

const sessionValueUserID = "user_id"

func (c *Controller) getSession(r *http.Request) (*sessions.Session, error) {
	return c.SessionStore.Get(r, isopodSession)
}

func (c *Controller) currentUserID(r *http.Request) (string, bool) {
	if c.SessionStore == nil {
		return "", false
	}
	s, err := c.getSession(r)
	if err != nil {
		return "", false
	}
	v, ok := s.Values[sessionValueUserID].(string)
	return v, ok && v != ""
}

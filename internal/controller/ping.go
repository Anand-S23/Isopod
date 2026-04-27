package controller

import (
	"net/http"
)

func (c *Controller) Ping(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, "Pong")
}

package controller

import (
	"context"
	"encoding/json"
	"net/http"
)

type Controller struct {
    Ctx          context.Context
    production   bool
}

func NewController(ctx context.Context, production bool) *Controller {
    return &Controller {
        Ctx: ctx,
        production: production,
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


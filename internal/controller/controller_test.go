package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Anand-S23/isopod/internal/config"
)

func TestNewController(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	env := &config.EnvVars{
		PRODUCTION:     true,
		BASE_URL:       "http://localhost:8080",
		CSRF_TOKEN:     "test_token",
		ENCRYPTION_KEY: []byte("test_key_123456"),
	}

	c := NewController(nil, env, ctx)
	if c.Ctx != ctx {
		t.Error("expected Ctx to be the passed context")
	}
	if !c.Production {
		t.Error("expected production true")
	}

	env2 := &config.EnvVars{
		PRODUCTION:     false,
		BASE_URL:       "http://localhost:8080",
		CSRF_TOKEN:     "test_token",
		ENCRYPTION_KEY: []byte("test_key_123456"),
	}

	c2 := NewController(nil, env2, context.Background())
	if c2.Production {
		t.Error("expected production false")
	}
}

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	payload := map[string]any{"ok": true, "n": 3}
	if err := WriteJSON(rr, http.StatusTeapot, payload); err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusTeapot)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type = %q, want json", ct)
	}

	var got map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if got["ok"] != true || got["n"] != float64(3) {
		t.Errorf("decoded body = %#v", got)
	}
}

func TestErrMsg(t *testing.T) {
	got := ErrMsg("something went wrong")
	if got["error"] != "something went wrong" {
		t.Errorf("ErrMsg() = %#v", got)
	}
}

func TestController_Ping(t *testing.T) {
	env := &config.EnvVars{
		PRODUCTION:     false,
		BASE_URL:       "http://localhost:8080",
		CSRF_TOKEN:     "test_token",
		ENCRYPTION_KEY: []byte("test_key_123456"),
	}
	c := NewController(nil, env, context.Background())
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)

	if err := c.Ping(rr, req); err != nil {
		t.Fatal(err)
	}

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type = %q, want json", ct)
	}

	var body string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body != "Pong" {
		t.Errorf("body = %q, want Pong", body)
	}
}

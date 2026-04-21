package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Anand-S23/isopod/internal/config"
	"github.com/Anand-S23/isopod/internal/controller"
)

func TestNewRouter_ping(t *testing.T) {
	env := &config.EnvVars{
		PRODUCTION:     false,
		BASE_URL:       "http://localhost:8080",
		CSRF_TOKEN:     "test_token",
		ENCRYPTION_KEY: []byte("test_key_123456"),
	}
	c := controller.NewController(nil, env, context.Background())
	mux := NewRouter(c)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	var body string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body != "Pong" {
		t.Errorf("body = %q, want Pong", body)
	}
}

func TestNewRouter_pingMethodNotMatched(t *testing.T) {
	env := &config.EnvVars{
		PRODUCTION:     false,
		BASE_URL:       "http://localhost:8080",
		CSRF_TOKEN:     "test_token",
		ENCRYPTION_KEY: []byte("test_key_123456"),
	}
	c := controller.NewController(nil, env, context.Background())
	mux := NewRouter(c)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /ping status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestFn_errorWritesJSON(t *testing.T) {
	h := Fn(func(http.ResponseWriter, *http.Request) error {
		return errors.New("handler failed")
	})

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type = %q", ct)
	}
	var got map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got["error"] != "handler failed" {
		t.Errorf("body = %#v", got)
	}
}

func TestFn_nilErrorNoExtraWrite(t *testing.T) {
	h := Fn(func(w http.ResponseWriter, _ *http.Request) error {
		w.WriteHeader(http.StatusNoContent)
		return nil
	})

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestNewCorsRouter_allowedOrigin(t *testing.T) {
	env := &config.EnvVars{
		PRODUCTION:     false,
		BASE_URL:       "http://localhost:8080",
		CSRF_TOKEN:     "test_token",
		ENCRYPTION_KEY: []byte("test_key_123456"),
	}
	c := controller.NewController(nil, env, context.Background())
	mux := NewRouter(c)
	h := NewCorsRouter(mux, "https://app.example")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("Origin", "https://app.example")
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d", rr.Code)
	}
	if aco := rr.Header().Get("Access-Control-Allow-Origin"); aco != "https://app.example" {
		t.Errorf("Access-Control-Allow-Origin = %q, want https://app.example", aco)
	}
}

func TestNewCorsRouter_preflight(t *testing.T) {
	env := &config.EnvVars{
		PRODUCTION:     false,
		BASE_URL:       "http://localhost:8080",
		CSRF_TOKEN:     "test_token",
		ENCRYPTION_KEY: []byte("test_key_123456"),
	}
	c := controller.NewController(nil, env, context.Background())
	mux := NewRouter(c)
	h := NewCorsRouter(mux, "https://app.example")

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/ping", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	h.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("expected Access-Control-Allow-Origin on preflight response")
	}
}

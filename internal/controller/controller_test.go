package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewController(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewController(ctx, true)
	if c.Ctx != ctx {
		t.Error("expected Ctx to be the passed context")
	}
	if !c.production {
		t.Error("expected production true")
	}

	c2 := NewController(context.Background(), false)
	if c2.production {
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
	c := NewController(context.Background(), false)
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


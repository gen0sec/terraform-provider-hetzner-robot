package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResetServer(t *testing.T) {
	var gotMethod, gotPath, gotType string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotType = r.PostFormValue("type")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"reset":{"server_ip":"1.2.3.4","type":["hw"]}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("user", "pass", srv.URL)
	if err := c.resetServer(context.Background(), 123456, "hw"); err != nil {
		t.Fatalf("resetServer returned error: %v", err)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("method = %s, want POST", gotMethod)
	}
	if gotPath != "/reset/123456" {
		t.Errorf("path = %s, want /reset/123456", gotPath)
	}
	if gotType != "hw" {
		t.Errorf("type = %s, want hw", gotType)
	}
}

func TestResetServer_acceptsAccepted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("user", "pass", srv.URL)
	if err := c.resetServer(context.Background(), 1, "power"); err != nil {
		t.Fatalf("resetServer should accept 202 Accepted, got: %v", err)
	}
}

func TestResetServer_errorStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":{"status":409,"code":"RESET_MANUAL_ACTIVE"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("user", "pass", srv.URL)
	if err := c.resetServer(context.Background(), 1, "hw"); err == nil {
		t.Fatal("expected error on 409, got nil")
	}
}

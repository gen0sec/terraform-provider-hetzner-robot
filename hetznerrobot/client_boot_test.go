package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteBootProfile(t *testing.T) {
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	if err := c.deleteBootProfile(context.Background(), 42, "rescue"); err != nil {
		t.Fatalf("deleteBootProfile error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
	if gotPath != "/boot/42/rescue" {
		t.Errorf("path = %s, want /boot/42/rescue", gotPath)
	}
}

func TestDeleteBootProfile_tolerates404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	if err := c.deleteBootProfile(context.Background(), 1, "linux"); err != nil {
		t.Fatalf("expected 404 to be tolerated, got: %v", err)
	}
}

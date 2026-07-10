package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestSetBootProfile_vnc(t *testing.T) {
	var gotPath string
	var gotForm url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotForm = r.PostForm
		_, _ = w.Write([]byte(`{"boot":{"vnc":{"active":true,"dist":"Ubuntu","lang":"en","server_number":42}}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	bp, err := c.setBootProfile(context.Background(), 42, "vnc", "Ubuntu", "en", nil)
	if err != nil {
		t.Fatalf("setBootProfile(vnc) error: %v", err)
	}
	if gotPath != "/boot/42/vnc" {
		t.Errorf("path = %s, want /boot/42/vnc", gotPath)
	}
	if gotForm.Get("dist") != "Ubuntu" || gotForm.Get("lang") != "en" {
		t.Errorf("form = %v, want dist=Ubuntu lang=en", gotForm)
	}
	if bp.ActiveProfile != "vnc" || bp.OperatingSystem != "Ubuntu" {
		t.Errorf("parsed = %+v", bp)
	}
}

func TestSetBootProfile_windows(t *testing.T) {
	var gotPath string
	var gotForm url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotForm = r.PostForm
		_, _ = w.Write([]byte(`{"boot":{"windows":{"active":true,"os":"standard","lang":"en_US","server_number":42}}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	bp, err := c.setBootProfile(context.Background(), 42, "windows", "standard", "en_US", nil)
	if err != nil {
		t.Fatalf("setBootProfile(windows) error: %v", err)
	}
	if gotPath != "/boot/42/windows" {
		t.Errorf("path = %s, want /boot/42/windows", gotPath)
	}
	if gotForm.Get("os") != "standard" || gotForm.Get("lang") != "en_US" {
		t.Errorf("form = %v, want os=standard lang=en_US", gotForm)
	}
	if bp.ActiveProfile != "windows" {
		t.Errorf("parsed = %+v", bp)
	}
}

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

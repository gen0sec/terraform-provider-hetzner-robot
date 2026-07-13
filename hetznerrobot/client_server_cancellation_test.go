package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/server/42/cancellation" {
			t.Errorf("path = %s, want /server/42/cancellation", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"cancellation":{"server_number":42,"server_ip":"1.2.3.4","cancellation_date":"2026-08-01","cancelled":false,"earliest_cancellation_date":"2026-08-01"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	cancel, err := c.getCancellation(context.Background(), 42)
	if err != nil {
		t.Fatalf("getCancellation error: %v", err)
	}
	if cancel.ServerNumber != 42 || cancel.Cancelled || cancel.EarliestCancellationDate != "2026-08-01" {
		t.Errorf("parsed cancellation = %+v", cancel)
	}
}

func TestCancelServer(t *testing.T) {
	var gotMethod, gotPath string
	var gotForm url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotForm = r.PostForm
		_, _ = w.Write([]byte(`{"cancellation":{"server_number":42,"cancellation_date":"now","cancelled":true}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	cancel, err := c.cancelServer(context.Background(), 42, "now", "no longer needed")
	if err != nil {
		t.Fatalf("cancelServer error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/server/42/cancellation" {
		t.Errorf("request = %s %s", gotMethod, gotPath)
	}
	if gotForm.Get("cancellation_date") != "now" || gotForm.Get("cancellation_reason") != "no longer needed" {
		t.Errorf("form = %v", gotForm)
	}
	if !cancel.Cancelled {
		t.Errorf("expected cancelled=true, got %+v", cancel)
	}
}

func TestRevokeCancellation(t *testing.T) {
	var gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	if err := c.revokeCancellation(context.Background(), 42); err != nil {
		t.Fatalf("revokeCancellation error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
}

func TestRenameServer(t *testing.T) {
	var gotMethod, gotPath, gotName string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotName = r.PostFormValue("server_name")
		_, _ = w.Write([]byte(`{"server":{"server_number":42,"server_name":"k8s-worker-1","server_ip":"1.2.3.4"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	s, err := c.renameServer(context.Background(), 42, "k8s-worker-1")
	if err != nil {
		t.Fatalf("renameServer error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/server/42" || gotName != "k8s-worker-1" {
		t.Errorf("request = %s %s name=%s", gotMethod, gotPath, gotName)
	}
	if s.ServerName != "k8s-worker-1" {
		t.Errorf("parsed server = %+v", s)
	}
}

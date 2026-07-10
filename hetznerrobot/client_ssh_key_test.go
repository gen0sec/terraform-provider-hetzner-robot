package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSshKeyByName(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/key" {
			t.Errorf("path = %s, want /key", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[
			{"key":{"name":"other","fingerprint":"11:11:11"}},
			{"key":{"name":"k8s-node-key","fingerprint":"aa:bb:cc","type":"ED25519","size":256}}
		]`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	key, err := c.getSshKeyByName(context.Background(), "k8s-node-key")
	if err != nil {
		t.Fatalf("getSshKeyByName error: %v", err)
	}
	if key.Fingerprint != "aa:bb:cc" || key.Type != "ED25519" || key.Size != 256 {
		t.Errorf("parsed key = %+v", key)
	}
}

func TestGetSshKeyByName_notFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"key":{"name":"other","fingerprint":"11:11:11"}}]`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	if _, err := c.getSshKeyByName(context.Background(), "missing"); err == nil {
		t.Fatal("expected error for missing key name, got nil")
	}
}

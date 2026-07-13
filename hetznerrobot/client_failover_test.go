package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFailover(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/failover/1.2.3.4" {
			t.Errorf("path = %s, want /failover/1.2.3.4", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"failover":{"ip":"1.2.3.4","netmask":"255.255.255.255","server_ip":"5.6.7.8","server_number":42,"active_server_ip":"5.6.7.8"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	fo, err := c.getFailover(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatalf("getFailover error: %v", err)
	}
	if fo.IP != "1.2.3.4" || fo.ServerNumber != 42 || fo.ActiveServerIP != "5.6.7.8" {
		t.Errorf("parsed failover = %+v", fo)
	}
}

func TestSetFailover(t *testing.T) {
	var gotMethod, gotPath, gotActive string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotActive = r.PostFormValue("active_server_ip")
		_, _ = w.Write([]byte(`{"failover":{"ip":"1.2.3.4","active_server_ip":"9.9.9.9","server_number":43}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	fo, err := c.setFailover(context.Background(), "1.2.3.4", "9.9.9.9")
	if err != nil {
		t.Fatalf("setFailover error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/failover/1.2.3.4" || gotActive != "9.9.9.9" {
		t.Errorf("request = %s %s active=%s", gotMethod, gotPath, gotActive)
	}
	if fo.ActiveServerIP != "9.9.9.9" {
		t.Errorf("parsed failover = %+v", fo)
	}
}

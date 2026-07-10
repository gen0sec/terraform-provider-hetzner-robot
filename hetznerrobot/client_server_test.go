package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetServer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/server/42" {
			t.Errorf("path = %s, want /server/42", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"server":{"server_ip":"1.2.3.4","server_number":42,"server_name":"node1","dc":"FSN1","status":"ready"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	s, err := c.getServer(context.Background(), 42)
	if err != nil {
		t.Fatalf("getServer error: %v", err)
	}
	if s.ServerNumber != 42 || s.ServerName != "node1" || s.ServerIP != "1.2.3.4" || s.DataCenter != "FSN1" {
		t.Errorf("parsed server = %+v", s)
	}
}

func TestGetServers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/server" {
			t.Errorf("path = %s, want /server", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[{"server":{"server_number":1,"server_name":"a"}},{"server":{"server_number":2,"server_name":"b"}}]`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	list, err := c.getServers(context.Background())
	if err != nil {
		t.Fatalf("getServers error: %v", err)
	}
	if len(list) != 2 || list[0].ServerNumber != 1 || list[1].ServerName != "b" {
		t.Errorf("parsed servers = %+v", list)
	}
}

package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRdns(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rdns/1.2.3.4" {
			t.Errorf("path = %s, want /rdns/1.2.3.4", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"rdns":{"ip":"1.2.3.4","ptr":"node1.example.com"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	rdns, err := c.getRdns(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatalf("getRdns error: %v", err)
	}
	if rdns.IP != "1.2.3.4" || rdns.PTR != "node1.example.com" {
		t.Errorf("parsed rdns = %+v", rdns)
	}
}

func TestSetRdns(t *testing.T) {
	var gotMethod, gotPath, gotPtr string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotPtr = r.PostFormValue("ptr")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"rdns":{"ip":"1.2.3.4","ptr":"node1.example.com"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	rdns, err := c.setRdns(context.Background(), "1.2.3.4", "node1.example.com")
	if err != nil {
		t.Fatalf("setRdns error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/rdns/1.2.3.4" || gotPtr != "node1.example.com" {
		t.Errorf("request = %s %s ptr=%s", gotMethod, gotPath, gotPtr)
	}
	if rdns.PTR != "node1.example.com" {
		t.Errorf("parsed rdns = %+v", rdns)
	}
}

func TestDeleteRdns(t *testing.T) {
	var gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	if err := c.deleteRdns(context.Background(), "1.2.3.4"); err != nil {
		t.Fatalf("deleteRdns error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
}

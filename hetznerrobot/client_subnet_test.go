package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSubnet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/subnet/2a01:4f8::" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"subnet":{"ip":"2a01:4f8::","mask":64,"gateway":"2a01:4f8::1","server_number":42,"failover":false,"traffic_warnings":true}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	sn, err := c.getSubnet(context.Background(), "2a01:4f8::")
	if err != nil {
		t.Fatalf("getSubnet error: %v", err)
	}
	if sn.Mask != 64 || sn.ServerNumber != 42 || !sn.TrafficWarnings {
		t.Errorf("parsed subnet = %+v", sn)
	}
}

func TestSetSubnet(t *testing.T) {
	var gotMethod, gotDaily string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_ = r.ParseForm()
		gotDaily = r.PostFormValue("traffic_daily")
		_, _ = w.Write([]byte(`{"subnet":{"ip":"2a01:4f8::","mask":64,"traffic_daily":50}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	sn, err := c.setSubnet(context.Background(), "2a01:4f8::", true, 10, 50, 100)
	if err != nil {
		t.Fatalf("setSubnet error: %v", err)
	}
	if gotMethod != http.MethodPost || gotDaily != "50" || sn.TrafficDaily != 50 {
		t.Errorf("method=%s daily=%s subnet=%+v", gotMethod, gotDaily, sn)
	}
}

func TestCreateSubnetMac(t *testing.T) {
	var gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_, _ = w.Write([]byte(`{"mac":{"ip":"2a01:4f8::","mac":"00:aa:bb:cc:dd:ee","possible":true}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	mac, err := c.createSubnetMac(context.Background(), "2a01:4f8::")
	if err != nil {
		t.Fatalf("createSubnetMac error: %v", err)
	}
	if gotMethod != http.MethodPut || mac.Mac != "00:aa:bb:cc:dd:ee" {
		t.Errorf("method=%s mac=%+v", gotMethod, mac)
	}
}

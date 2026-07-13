package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ip/1.2.3.4" {
			t.Errorf("path = %s, want /ip/1.2.3.4", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ip":{"ip":"1.2.3.4","server_ip":"1.2.3.4","server_number":42,"locked":false,"traffic_warnings":true,"traffic_hourly":100,"traffic_daily":200,"traffic_monthly":300}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	ip, err := c.getIP(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatalf("getIP error: %v", err)
	}
	if ip.ServerNumber != 42 || !ip.TrafficWarnings || ip.TrafficMonthly != 300 {
		t.Errorf("parsed ip = %+v", ip)
	}
}

func TestSetIP(t *testing.T) {
	var gotMethod, gotHourly string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_ = r.ParseForm()
		gotHourly = r.PostFormValue("traffic_hourly")
		_, _ = w.Write([]byte(`{"ip":{"ip":"1.2.3.4","traffic_hourly":100}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	ip, err := c.setIP(context.Background(), "1.2.3.4", true, 100, 200, 300)
	if err != nil {
		t.Fatalf("setIP error: %v", err)
	}
	if gotMethod != http.MethodPost || gotHourly != "100" {
		t.Errorf("request = %s hourly=%s", gotMethod, gotHourly)
	}
	if ip.TrafficHourly != 100 {
		t.Errorf("parsed ip = %+v", ip)
	}
}

func TestCreateIPMac(t *testing.T) {
	var gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		if r.URL.Path != "/ip/1.2.3.4/mac" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"mac":{"ip":"1.2.3.4","mac":"00:11:22:33:44:55","possible":true}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	mac, err := c.createIPMac(context.Background(), "1.2.3.4")
	if err != nil {
		t.Fatalf("createIPMac error: %v", err)
	}
	if gotMethod != http.MethodPut || mac.Mac != "00:11:22:33:44:55" {
		t.Errorf("method=%s mac=%+v", gotMethod, mac)
	}
}

func TestDeleteIPMac(t *testing.T) {
	var gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	if err := c.deleteIPMac(context.Background(), "1.2.3.4"); err != nil {
		t.Fatalf("deleteIPMac error: %v", err)
	}
	if gotMethod != http.MethodDelete {
		t.Errorf("method = %s, want DELETE", gotMethod)
	}
}

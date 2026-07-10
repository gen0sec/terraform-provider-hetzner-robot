package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetServerAddonProducts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/order/server_addon/2093885/product" {
			t.Errorf("path = %s, want /order/server_addon/2093885/product", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[
			{"product":{"id":"failover_ip","name":"Additional single IP","type":"ip","price":{"price":{"net":"3.36","gross":"4.00"}}}}
		]`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	addons, err := c.getServerAddonProducts(context.Background(), 2093885)
	if err != nil {
		t.Fatalf("getServerAddonProducts error: %v", err)
	}
	if len(addons) != 1 {
		t.Fatalf("got %d addons, want 1", len(addons))
	}
	a := addons[0]
	if a.ID != "failover_ip" || a.Type != "ip" || a.PriceNet != "3.36" || a.PriceGross != "4.00" {
		t.Errorf("parsed addon = %+v", a)
	}
}

func TestCreateServerAddonOrder(t *testing.T) {
	var gotPath string
	var gotForm url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotForm = r.PostForm
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"transaction":{"id":"B-A1","status":"in process","server_number":2093885}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	tx, err := c.createServerAddonOrder(context.Background(), 2093885, "failover_ip", true)
	if err != nil {
		t.Fatalf("createServerAddonOrder error: %v", err)
	}
	if gotPath != "/order/server_addon/transaction" {
		t.Errorf("path = %s", gotPath)
	}
	if gotForm.Get("server_number") != "2093885" || gotForm.Get("product_id") != "failover_ip" || gotForm.Get("test") != "true" {
		t.Errorf("form = %v", gotForm)
	}
	if tx.ID != "B-A1" || tx.ServerNumber != 2093885 {
		t.Errorf("parsed transaction = %+v", tx)
	}
}

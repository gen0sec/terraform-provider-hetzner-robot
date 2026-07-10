package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetServerMarketProducts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/order/server_market/product" {
			t.Errorf("path = %s, want /order/server_market/product", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[
			{"product":{"id":"1234","name":"SB with Ryzen","description":["AMD Ryzen 7"],"traffic":"unlimited","dist":["Rescue system"],"lang":["en"],"arch":["64"],"cpu":"AMD Ryzen 7 1700X","cpu_benchmark":8877,"memory_size":64,"hdd_size":3000,"hdd_text":"2x SSD","hdd_count":2,"datacenter":"FSN1-DC5","network_speed":"1 Gbit","price":"40.00","fixed_price":true,"next_reduce":0}}
		]`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	products, err := c.getServerMarketProducts(context.Background())
	if err != nil {
		t.Fatalf("getServerMarketProducts error: %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("got %d products, want 1", len(products))
	}
	p := products[0]
	if p.ID != "1234" || p.MemorySize != 64 || p.CPUBenchmark != 8877 || p.Datacenter != "FSN1-DC5" || !p.FixedPrice {
		t.Errorf("parsed market product = %+v", p)
	}
}

func TestCreateServerMarketOrder(t *testing.T) {
	var gotPath string
	var gotForm url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotForm = r.PostForm
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"transaction":{"id":"B-9","status":"in process"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	tx, err := c.createServerMarketOrder(context.Background(), HetznerRobotServerOrderRequest{
		ProductID:      "1234",
		Location:       "SHOULD-BE-IGNORED",
		AuthorizedKeys: []string{"aa:bb"},
		Test:           true,
	})
	if err != nil {
		t.Fatalf("createServerMarketOrder error: %v", err)
	}
	if gotPath != "/order/server_market/transaction" {
		t.Errorf("path = %s", gotPath)
	}
	if gotForm.Get("product_id") != "1234" || gotForm.Get("test") != "true" {
		t.Errorf("form = %v", gotForm)
	}
	if gotForm.Has("location") {
		t.Errorf("market order must not send location, got %q", gotForm.Get("location"))
	}
	if tx.ID != "B-9" {
		t.Errorf("parsed transaction = %+v", tx)
	}
}

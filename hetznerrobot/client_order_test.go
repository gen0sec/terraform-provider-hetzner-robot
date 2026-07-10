package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetServerProducts(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/order/server/product" {
			t.Errorf("path = %s, want /order/server/product", r.URL.Path)
		}
		_, _ = w.Write([]byte(`[
			{"product":{"id":"EX44","name":"Dedicated Root Server EX44","description":["Intel Core i5-13500","64 GB DDR4"],"traffic":"unlimited","location":["FSN1","NBG1"],"prices":[{"location":"FSN1","price":{"net":"39.00","gross":"46.41"},"price_setup":{"net":"39.00","gross":"46.41"}}]}}
		]`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	products, err := c.getServerProducts(context.Background())
	if err != nil {
		t.Fatalf("getServerProducts error: %v", err)
	}
	if len(products) != 1 {
		t.Fatalf("got %d products, want 1", len(products))
	}
	p := products[0]
	if p.ID != "EX44" || p.Traffic != "unlimited" || len(p.Description) != 2 || len(p.Locations) != 2 {
		t.Errorf("parsed product = %+v", p)
	}
	if len(p.Prices) != 1 || p.Prices[0].Location != "FSN1" || p.Prices[0].MonthlyNet != "39.00" || p.Prices[0].SetupGross != "46.41" {
		t.Errorf("parsed prices = %+v", p.Prices)
	}
}

func TestCreateServerOrder_test(t *testing.T) {
	var gotPath string
	var gotForm url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_ = r.ParseForm()
		gotForm = r.PostForm
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"transaction":{"id":"B-123","status":"in process","server_number":null,"server_ip":null}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	tx, err := c.createServerOrder(context.Background(), HetznerRobotServerOrderRequest{
		ProductID:      "EX44",
		Location:       "FSN1",
		Dist:           "Ubuntu 24.04 LTS minimal",
		AuthorizedKeys: []string{"aa:bb", "cc:dd"},
		Test:           true,
	})
	if err != nil {
		t.Fatalf("createServerOrder error: %v", err)
	}
	if gotPath != "/order/server/transaction" {
		t.Errorf("path = %s, want /order/server/transaction", gotPath)
	}
	if gotForm.Get("product_id") != "EX44" || gotForm.Get("location") != "FSN1" || gotForm.Get("test") != "true" {
		t.Errorf("form = %v", gotForm)
	}
	if keys := gotForm["authorized_key[]"]; len(keys) != 2 {
		t.Errorf("authorized_key[] = %v, want 2 entries", keys)
	}
	if tx.ID != "B-123" || tx.Status != "in process" || tx.ServerNumber != 0 {
		t.Errorf("parsed transaction = %+v", tx)
	}
}

func TestGetServerOrder(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/order/server/transaction/B-123" {
			t.Errorf("path = %s, want /order/server/transaction/B-123", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"transaction":{"id":"B-123","status":"ready","server_number":2093885,"server_ip":"1.2.3.4"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	tx, err := c.getServerOrder(context.Background(), "B-123")
	if err != nil {
		t.Fatalf("getServerOrder error: %v", err)
	}
	if tx.Status != "ready" || tx.ServerNumber != 2093885 || tx.ServerIP != "1.2.3.4" {
		t.Errorf("parsed transaction = %+v", tx)
	}
}

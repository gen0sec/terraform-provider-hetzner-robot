package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourceSubnet_CreateAndRead(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"subnet":{"ip":"2a01:4f8::","mask":64,"gateway":"2a01:4f8::1","server_number":42,"failover":false,"traffic_warnings":true,"traffic_monthly":50}}`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceSubnet()
	d := res.TestResourceData()
	d.Set("subnet_ip", "2a01:4f8::")
	d.Set("traffic_warnings", true)
	d.Set("traffic_monthly", 50)

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if d.Id() != "2a01:4f8::" || d.Get("mask").(int) != 64 || d.Get("gateway").(string) != "2a01:4f8::1" {
		t.Errorf("subnet = id:%s mask:%v gw:%v", d.Id(), d.Get("mask"), d.Get("gateway"))
	}
	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
}

func TestResourceSubnetMac_CreateReadDelete(t *testing.T) {
	var deleted bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodGet:
			_, _ = w.Write([]byte(`{"mac":{"ip":"2a01:4f8::","mac":"00:aa:bb:cc:dd:ee","possible":true}}`))
		case http.MethodDelete:
			deleted = true
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceSubnetMac()
	d := res.TestResourceData()
	d.Set("subnet_ip", "2a01:4f8::")

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if d.Get("mac").(string) != "00:aa:bb:cc:dd:ee" {
		t.Errorf("mac = %v", d.Get("mac"))
	}
	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if diags := res.DeleteContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("delete diags: %v", diags)
	}
	if !deleted || d.Id() != "" {
		t.Errorf("delete state: deleted=%v id=%q", deleted, d.Id())
	}
}

func TestDataSubnet_Read(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"subnet":{"ip":"2a01:4f8::","mask":56,"failover":true,"locked":false}}`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	ds := dataSubnet()
	d := ds.TestResourceData()
	d.Set("subnet_ip", "2a01:4f8::")
	if diags := ds.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if d.Id() != "2a01:4f8::" || d.Get("mask").(int) != 56 || !d.Get("failover").(bool) {
		t.Errorf("subnet data = id:%s mask:%v failover:%v", d.Id(), d.Get("mask"), d.Get("failover"))
	}
}

func TestDataTraffic_Read(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"traffic":{"data":{"1.2.3.4":{"in":"1.5","out":"2.5","sum":"4.0"}}}}`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	ds := dataTraffic()
	d := ds.TestResourceData()
	d.Set("ip", "1.2.3.4")
	d.Set("type", "month")
	d.Set("from", "2024-01")
	d.Set("to", "2024-12")
	if diags := ds.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if d.Get("in").(string) != "1.5" || d.Get("sum").(string) != "4.0" {
		t.Errorf("traffic = in:%v sum:%v", d.Get("in"), d.Get("sum"))
	}
	if d.Id() == "" {
		t.Errorf("traffic data id not set")
	}
}

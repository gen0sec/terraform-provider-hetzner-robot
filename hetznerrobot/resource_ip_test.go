package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourceIP_CreateAndRead(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ip/1.2.3.4" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"ip":{"ip":"1.2.3.4","server_ip":"1.2.3.4","server_number":42,"separate_mac":"aa:bb:cc:dd:ee:ff","traffic_warnings":true,"traffic_hourly":100,"traffic_daily":200,"traffic_monthly":300}}`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceIP()
	d := res.TestResourceData()
	d.Set("ip", "1.2.3.4")
	d.Set("traffic_warnings", true)
	d.Set("traffic_hourly", 100)
	d.Set("traffic_daily", 200)
	d.Set("traffic_monthly", 300)

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if d.Id() != "1.2.3.4" {
		t.Errorf("id = %q, want 1.2.3.4", d.Id())
	}
	if d.Get("server_number").(int) != 42 {
		t.Errorf("server_number = %v", d.Get("server_number"))
	}
	if d.Get("separate_mac").(string) != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("separate_mac = %v", d.Get("separate_mac"))
	}

	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if !d.Get("traffic_warnings").(bool) {
		t.Errorf("traffic_warnings not read back")
	}
}

func TestResourceIP_DeleteIsNoOp(t *testing.T) {
	// Delete must not call the API and simply clears the ID.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("delete unexpectedly called the API: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceIP()
	d := res.TestResourceData()
	d.SetId("1.2.3.4")
	if diags := res.DeleteContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("delete diags: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("id after delete = %q, want empty", d.Id())
	}
}

func TestResourceIPMac_CreateReadDelete(t *testing.T) {
	var deleted bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodGet:
			_, _ = w.Write([]byte(`{"mac":{"ip":"1.2.3.4","mac":"00:11:22:33:44:55","possible":true}}`))
		case http.MethodDelete:
			deleted = true
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceIPMac()
	d := res.TestResourceData()
	d.Set("ip", "1.2.3.4")

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if d.Get("mac").(string) != "00:11:22:33:44:55" {
		t.Errorf("mac = %v", d.Get("mac"))
	}
	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if diags := res.DeleteContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("delete diags: %v", diags)
	}
	if !deleted {
		t.Errorf("delete did not hit the API")
	}
	if d.Id() != "" {
		t.Errorf("id after delete = %q", d.Id())
	}
}

func TestDataIP_Read(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ip":{"ip":"1.2.3.4","server_number":7,"locked":true,"traffic_monthly":42}}`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	ds := dataIP()
	d := ds.TestResourceData()
	d.Set("ip", "1.2.3.4")
	if diags := ds.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if d.Id() != "1.2.3.4" || d.Get("server_number").(int) != 7 || !d.Get("locked").(bool) {
		t.Errorf("data ip = id:%s sn:%v locked:%v", d.Id(), d.Get("server_number"), d.Get("locked"))
	}
}

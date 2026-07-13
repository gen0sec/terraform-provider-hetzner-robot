package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourceStorageBox_CreateAndRead(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storagebox/123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		// Both the POST (update) and GET (read) return the box.
		_, _ = w.Write([]byte(`{"storagebox":{"id":123,"login":"u123","name":"k8s-backups","product":"BX11","location":"FSN1","disk_quota":1024,"ssh":true,"samba":false,"webdav":false,"external_reachability":true,"zfs":false}}`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceStorageBox()
	d := res.TestResourceData()
	d.Set("storagebox_id", 123)
	d.Set("name", "k8s-backups")
	d.Set("ssh", true)
	d.Set("external_reachability", true)

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if d.Id() != "123" {
		t.Errorf("id = %q, want 123", d.Id())
	}
	if !d.Get("ssh").(bool) || d.Get("product").(string) != "BX11" || d.Get("location").(string) != "FSN1" {
		t.Errorf("box = ssh:%v product:%v location:%v", d.Get("ssh"), d.Get("product"), d.Get("location"))
	}

	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if d.Get("name").(string) != "k8s-backups" {
		t.Errorf("name after read = %v", d.Get("name"))
	}
}

func TestResourceStorageBox_DeleteIsNoOp(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("delete unexpectedly hit the API: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceStorageBox()
	d := res.TestResourceData()
	d.SetId("123")
	if diags := res.DeleteContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("delete diags: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("id after delete = %q", d.Id())
	}
}

func TestResourceStorageBoxSnapshot_CreateReadDelete(t *testing.T) {
	var deleted bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost:
			_, _ = w.Write([]byte(`{"snapshot":{"name":"2024-01-01T12-00-00","timestamp":"2024-01-01T12:00:00+00:00","size":10,"automatic":false}}`))
		case r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`[{"snapshot":{"name":"2024-01-01T12-00-00","timestamp":"2024-01-01T12:00:00+00:00","size":10,"automatic":false}}]`))
		case r.Method == http.MethodDelete:
			deleted = true
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceStorageBoxSnapshot()
	d := res.TestResourceData()
	d.Set("storagebox_id", 123)

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if d.Get("name").(string) != "2024-01-01T12-00-00" || d.Get("size").(int) != 10 {
		t.Errorf("snapshot = name:%v size:%v", d.Get("name"), d.Get("size"))
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

func TestResourceStorageBoxSnapshot_ReadGoneClearsState(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// snapshot no longer present in the list
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceStorageBoxSnapshot()
	d := res.TestResourceData()
	d.Set("storagebox_id", 123)
	d.Set("name", "2024-01-01T12-00-00")
	d.SetId("123/2024-01-01T12-00-00")
	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("expected state cleared, id = %q", d.Id())
	}
}

func TestResourceStorageBoxSubaccount_CreateReadUpdateDelete(t *testing.T) {
	var putCalled, deleteCalled bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			_, _ = w.Write([]byte(`{"subaccount":{"username":"u123-sub1","password":"secret","homedirectory":"/backups","ssh":true,"server":"u123.your-storagebox.de"}}`))
		case http.MethodGet:
			_, _ = w.Write([]byte(`[{"subaccount":{"username":"u123-sub1","homedirectory":"/backups","ssh":true,"readonly":true,"server":"u123.your-storagebox.de"}}]`))
		case http.MethodPut:
			putCalled = true
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			deleteCalled = true
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceStorageBoxSubaccount()
	d := res.TestResourceData()
	d.Set("storagebox_id", 123)
	d.Set("homedirectory", "/backups")
	d.Set("ssh", true)

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if d.Get("username").(string) != "u123-sub1" || d.Get("password").(string) != "secret" {
		t.Errorf("subaccount = user:%v pass:%v", d.Get("username"), d.Get("password"))
	}
	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if !d.Get("readonly").(bool) {
		t.Errorf("readonly not read back")
	}
	if diags := res.UpdateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("update diags: %v", diags)
	}
	if !putCalled {
		t.Errorf("update did not issue PUT")
	}
	if diags := res.DeleteContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("delete diags: %v", diags)
	}
	if !deleteCalled || d.Id() != "" {
		t.Errorf("delete state: called=%v id=%q", deleteCalled, d.Id())
	}
}

func TestDataStorageBox_Read(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storagebox/123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"storagebox":{"id":123,"name":"backup","product":"BX11","disk_quota":1024,"ssh":true}}`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	ds := dataStorageBox()
	d := ds.TestResourceData()
	d.Set("storagebox_id", 123)
	if diags := ds.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if d.Id() != "123" || d.Get("name").(string) != "backup" || d.Get("disk_quota").(int) != 1024 {
		t.Errorf("box data = id:%s name:%v quota:%v", d.Id(), d.Get("name"), d.Get("disk_quota"))
	}
}

func TestDataStorageBoxes_Read(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"storagebox":{"id":1,"name":"a"}},{"storagebox":{"id":2,"name":"b"}}]`))
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	ds := dataStorageBoxes()
	d := ds.TestResourceData()
	if diags := ds.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	boxes := d.Get("storageboxes").([]interface{})
	if len(boxes) != 2 {
		t.Fatalf("got %d boxes, want 2", len(boxes))
	}
	first := boxes[0].(map[string]interface{})
	if first["name"].(string) != "a" || first["storagebox_id"].(int) != 1 {
		t.Errorf("first box = %+v", first)
	}
}

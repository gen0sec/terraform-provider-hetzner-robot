package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStorageBox(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storagebox/123" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"storagebox":{"id":123,"login":"u123","name":"backup","product":"BX11","disk_quota":1024,"ssh":true,"samba":false}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	b, err := c.getStorageBox(context.Background(), 123)
	if err != nil {
		t.Fatalf("getStorageBox error: %v", err)
	}
	if b.ID != 123 || b.Name != "backup" || !b.SSH || b.DiskQuota != 1024 {
		t.Errorf("parsed box = %+v", b)
	}
}

func TestGetStorageBoxes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"storagebox":{"id":1,"name":"a"}},{"storagebox":{"id":2,"name":"b"}}]`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	boxes, err := c.getStorageBoxes(context.Background())
	if err != nil {
		t.Fatalf("getStorageBoxes error: %v", err)
	}
	if len(boxes) != 2 || boxes[0].ID != 1 || boxes[1].Name != "b" {
		t.Errorf("parsed boxes = %+v", boxes)
	}
}

func TestUpdateStorageBox(t *testing.T) {
	var gotMethod, gotName, gotSSH string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		_ = r.ParseForm()
		gotName = r.PostFormValue("storagebox_name")
		gotSSH = r.PostFormValue("ssh")
		_, _ = w.Write([]byte(`{"storagebox":{"id":123,"name":"renamed","ssh":true}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	b, err := c.updateStorageBox(context.Background(), 123, "renamed", true, false, false, false, false)
	if err != nil {
		t.Fatalf("updateStorageBox error: %v", err)
	}
	if gotMethod != http.MethodPost || gotName != "renamed" || gotSSH != "true" || b.Name != "renamed" {
		t.Errorf("request=%s name=%s ssh=%s box=%+v", gotMethod, gotName, gotSSH, b)
	}
}

func TestCreateStorageBoxSnapshot(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/storagebox/123/snapshot" || r.Method != http.MethodPost {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"snapshot":{"name":"2024-01-01T12-00-00","timestamp":"2024-01-01T12:00:00+00:00","size":10,"automatic":false}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	snap, err := c.createStorageBoxSnapshot(context.Background(), 123)
	if err != nil {
		t.Fatalf("createStorageBoxSnapshot error: %v", err)
	}
	if snap.Name != "2024-01-01T12-00-00" || snap.Size != 10 {
		t.Errorf("parsed snapshot = %+v", snap)
	}
}

func TestCreateStorageBoxSubaccount(t *testing.T) {
	var gotHomedir string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		gotHomedir = r.PostFormValue("homedirectory")
		_, _ = w.Write([]byte(`{"subaccount":{"username":"u123-sub1","password":"secret","homedirectory":"/backups","ssh":true}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	sub, err := c.createStorageBoxSubaccount(context.Background(), 123, "/backups", false, true, false, false, false, "note")
	if err != nil {
		t.Fatalf("createStorageBoxSubaccount error: %v", err)
	}
	if gotHomedir != "/backups" || sub.Username != "u123-sub1" || sub.Password != "secret" {
		t.Errorf("homedir=%s sub=%+v", gotHomedir, sub)
	}
}

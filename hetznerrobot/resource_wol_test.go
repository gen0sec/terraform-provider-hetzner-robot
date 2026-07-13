package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResourceWOL_CreateSendsPacket(t *testing.T) {
	var gotPath, gotMethod string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()
	c := NewHetznerRobotClient("u", "p", srv.URL)

	res := resourceWOL()
	d := res.TestResourceData()
	d.Set("server_number", 123456)

	if diags := res.CreateContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("create diags: %v", diags)
	}
	if gotMethod != http.MethodPost || gotPath != "/wol/123456" {
		t.Errorf("request = %s %s", gotMethod, gotPath)
	}
	if d.Id() != "123456" {
		t.Errorf("id = %q, want 123456", d.Id())
	}

	// Read is a no-op; Delete only clears state.
	if diags := res.ReadContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("read diags: %v", diags)
	}
	if diags := res.DeleteContext(context.Background(), d, c); diags.HasError() {
		t.Fatalf("delete diags: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("id after delete = %q", d.Id())
	}
}

package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendWOL(t *testing.T) {
	var gotMethod, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	if err := c.sendWOL(context.Background(), 42); err != nil {
		t.Fatalf("sendWOL error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/wol/42" {
		t.Errorf("request = %s %s", gotMethod, gotPath)
	}
}

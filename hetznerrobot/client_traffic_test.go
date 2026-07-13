package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTraffic(t *testing.T) {
	var gotType, gotIP string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/traffic" {
			t.Errorf("path = %s, want /traffic", r.URL.Path)
		}
		_ = r.ParseForm()
		gotType = r.PostFormValue("type")
		gotIP = r.PostFormValue("ip[]")
		_, _ = w.Write([]byte(`{"traffic":{"type":"day","from":"2024-01-01","to":"2024-01-02","data":{"1.2.3.4":{"in":"1.5","out":"2.5","sum":"4.0"}}}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("u", "p", srv.URL)
	tr, err := c.getTraffic(context.Background(), "1.2.3.4", "day", "2024-01-01", "2024-01-02")
	if err != nil {
		t.Fatalf("getTraffic error: %v", err)
	}
	if gotType != "day" || gotIP != "1.2.3.4" {
		t.Errorf("request type=%s ip=%s", gotType, gotIP)
	}
	if tr.In != "1.5" || tr.Out != "2.5" || tr.Sum != "4.0" {
		t.Errorf("parsed traffic = %+v", tr)
	}
}

package hetznerrobot

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCodeIsInExpected(t *testing.T) {
	cases := []struct {
		code     int
		expected []int
		want     bool
	}{
		{200, []int{200, 202}, true},
		{202, []int{200, 202}, true},
		{500, []int{200, 202}, false},
		{200, nil, false},
	}
	for _, tc := range cases {
		if got := codeIsInExpected(tc.code, tc.expected); got != tc.want {
			t.Errorf("codeIsInExpected(%d, %v) = %v, want %v", tc.code, tc.expected, got, tc.want)
		}
	}
}

func TestMakeAPICall_basicAuthAndBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if u, p, ok := r.BasicAuth(); !ok || u != "user" || p != "pass" {
			t.Errorf("basic auth = %q/%q ok=%v, want user/pass", u, p, ok)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("user", "pass", srv.URL)

	body, err := c.makeAPICall(context.Background(), http.MethodGet, srv.URL+"/x", nil, []int{http.StatusOK})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "ok" {
		t.Errorf("body = %q, want %q", body, "ok")
	}
}

func TestMakeAPICall_unexpectedStatusIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":{"status":500,"message":"boom"}}`))
	}))
	defer srv.Close()

	c := NewHetznerRobotClient("user", "pass", srv.URL)
	if _, err := c.makeAPICall(context.Background(), http.MethodGet, srv.URL+"/x", nil, []int{http.StatusOK}); err == nil {
		t.Fatal("expected error for unmatched status code, got nil")
	}
}

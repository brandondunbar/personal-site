// cmd/web/main_test.go
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootServesGreeting(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	routes().ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("want status 200 OK, got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading body: %v", err)
	}
	want := "Hello, personal-site is running!"
	if string(body) != want {
		t.Fatalf("unexpected body:\nwant: %q\ngot:  %q", want, string(body))
	}
}

func TestHealthz_OK(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	routes().ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("want 200 OK, got %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	if string(b) != "OK" {
		t.Fatalf("want body %q, got %q", "OK", string(b))
	}
	if ct := res.Header.Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("want Content-Type %q, got %q", "text/plain; charset=utf-8", ct)
	}
}

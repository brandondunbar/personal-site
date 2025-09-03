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

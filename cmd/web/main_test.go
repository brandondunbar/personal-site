// cmd/web/main_test.go
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
		t.Fatalf("want 200 OK, got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Hello from Go templates!") {
		t.Fatalf("want greeting in HTML, got: %s", string(body))
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

func TestRootRendersHTML(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	routes().ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("want status 200 OK, got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Fatalf("want Content-Type %q, got %q", "text/html; charset=utf-8", ct)
	}
	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "<h1>") {
		t.Fatalf("want HTML body containing <h1>, got: %s", string(body))
	}
}

func TestStaticServesFile(t *testing.T) {
	t.Parallel()

	// Arrange: create a temporary static directory with one file.
	tmp := t.TempDir()
	js := []byte(`console.log("ok");`)
	if err := os.WriteFile(filepath.Join(tmp, "app.js"), js, 0o644); err != nil {
		t.Fatalf("write tmp static: %v", err)
	}

	// Swap the static filesystem to point to our temp dir for this test.
	orig := staticFS
	staticFS = http.Dir(tmp)
	t.Cleanup(func() { staticFS = orig })

	// Act
	req := httptest.NewRequest(http.MethodGet, "/static/app.js", nil)
	rec := httptest.NewRecorder()
	routes().ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	// Assert
	if res.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(res.Body)
		t.Fatalf("want 200, got %d body=%s", res.StatusCode, string(b))
	}
	cc := res.Header.Get("Cache-Control")
	if cc == "" || !strings.Contains(cc, "max-age=") {
		t.Fatalf("expected Cache-Control with max-age, got %q", cc)
	}
	// Last-Modified should be set by http.FileServer
	if lm := res.Header.Get("Last-Modified"); lm == "" {
		t.Fatalf("expected Last-Modified header")
	}
}

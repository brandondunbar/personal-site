// cmd/web/main_test.go
package main

import (
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/brandondunbar/personal-site/internal/config"
)

func TestHealthz_OK(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/healthz")
	if err != nil {
		t.Fatalf("GET /healthz: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/plain") {
		t.Fatalf("Content-Type = %q, want text/plain", ct)
	}
}

func TestStaticServesFile_WithCacheControl(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	// File created by mustTestApp at /static/css/site.css
	resp, err := http.Get(srv.URL + "/static/css/site.css")
	if err != nil {
		t.Fatalf("GET /static/css/site.css: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	gotCC := resp.Header.Get("Cache-Control")
	wantCC := "public, max-age=31536000, immutable"
	if gotCC != wantCC {
		t.Fatalf("Cache-Control = %q, want %q", gotCC, wantCC)
	}
}

func TestHome_Renders_WithConfigAndYear(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	body, _ := ioReadAll(resp.Body)
	if !strings.Contains(body, "Home - Elliot Alderson") {
		t.Fatalf("body missing title with Site.Name: %q", body)
	}
	if !strings.Contains(body, "name@domain.com") {
		t.Fatalf("body missing Site.Email: %q", body)
	}
	year := strconv.Itoa(time.Now().Year())
	if !strings.Contains(body, year) {
		t.Fatalf("body missing current year %s: %q", year, body)
	}
}

/************ helpers ************/

func mustTestApp(t *testing.T) *App {
	t.Helper()

	// --- minimal templates exercising .Site and .Year ---
	const baseTpl = `{{define "base"}}<html><head><title>{{block "title" .}}x{{end}}</title></head><body>{{block "content" .}}{{end}}</body></html>{{end}}`
	const homeTpl = `{{define "title"}}Home - {{.Site.Name}}{{end}}{{define "content"}}Hello {{.Site.Email}} — {{.Year}}{{end}}`

	tpls := template.Must(template.New("base").Parse(baseTpl))
	template.Must(tpls.Parse(homeTpl))

	// --- temp static dir with one file at css/site.css ---
	td := t.TempDir()
	staticRoot := filepath.Join(td, "css")
	if err := os.MkdirAll(staticRoot, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(staticRoot, "site.css"), []byte("/* test */"), 0o644); err != nil {
		t.Fatalf("write css: %v", err)
	}

	// Build the app
	return &App{
		tpls:     tpls,
		staticFS: http.Dir(td), // so /static/css/site.css maps into td/css/site.css
		cfg: config.Config{
			Name:    "Elliot Alderson",
			Email:   "name@domain.com",
			Brand:   "mr.robot",
			Tagline: "Computer Repair with a Smile",
		},
	}
}

func ioReadAll(r io.Reader) (string, error) {
	var sb strings.Builder
	_, err := io.Copy(&sb, r)
	return sb.String(), err
}

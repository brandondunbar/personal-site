// cmd/web/main_test.go
package main

import (
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/brandondunbar/personal-site/internal/config"
)

func TestHealthz_OK(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.Routes())
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
	srv := httptest.NewServer(app.Routes())
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
	srv := httptest.NewServer(app.Routes())
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

func TestPanicRecovery_Returns500(t *testing.T) {
	app := mustTestApp(t)

	// A handler that panics before writing anything
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	srv := httptest.NewServer(app.recoverMiddleware(panicHandler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/anything")
	if err != nil {
		t.Fatalf("GET panic route: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", ct)
	}
}

// --- additional tests ---

// 1) Unknown routes return 404 (ensures mux wiring isn't shadowed)
func TestNotFound_Returns404(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/nope")
	if err != nil {
		t.Fatalf("GET /nope: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
}

// 2) HEAD on /healthz should also be OK
func TestHealthz_HEAD_OK(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.Routes())
	defer srv.Close()

	req, _ := http.NewRequest("HEAD", srv.URL+"/healthz", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("HEAD /healthz: %v", err)
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

// 3) Home sets HTML content-type (sanity on successful render path)
func TestHome_ContentTypeHTML(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", ct)
	}
}

// 4) If the template is broken, handler should 500 (exercise error branch)
func TestHome_TemplateError_Returns500Plain(t *testing.T) {
	app := mustTestApp(t)

	// Replace templates with one that lacks "base" to force ExecuteTemplate error
	tpls := template.Must(template.New("x").Parse(`{{define "notbase"}}noop{{end}}`))
	app.tpls = tpls

	srv := httptest.NewServer(app.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/plain") {
		t.Fatalf("Content-Type = %q, want text/plain", ct)
	}
	body, _ := ioReadAll(resp.Body)
	if !strings.Contains(body, "template error:") {
		t.Fatalf("body = %q, want contains 'template error:'", body)
	}
}

// 5) Recovery middleware renders friendly 500 page when a handler panics (using template if available)
func TestPanicRecovery_UsesErrorTemplateWhenAvailable(t *testing.T) {
	app := mustTestApp(t)

	// Provide a tiny error template so the middleware path is deterministic
	errTpl := template.Must(template.New("err").Parse(`{{define "error500"}}ERR {{.Year}} {{.Site.Name}}{{end}}`))
	app.tpls = errTpl

	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { panic("boom") })
	srv := httptest.NewServer(app.recoverMiddleware(panicHandler))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/panic")
	if err != nil {
		t.Fatalf("GET /panic: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", ct)
	}
	body, _ := ioReadAll(resp.Body)
	if !strings.Contains(body, "ERR") {
		t.Fatalf("error template not used, body=%q", body)
	}
}

// 6) Light concurrency probe on home route (helps catch data races with -race)
func TestHome_Concurrent_NoPanic(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.Routes())
	defer srv.Close()

	var wg sync.WaitGroup
	const N = 32
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			resp, err := http.Get(srv.URL + "/")
			if err != nil {
				t.Errorf("GET /: %v", err)
				return
			}
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()
	}
	wg.Wait()
}

// 7) Request id header present
func TestRequestID_HeaderPresent(t *testing.T) {
	app := mustTestApp(t)
	srv := httptest.NewServer(app.Routes())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()

	id := resp.Header.Get("X-Request-Id")
	if id == "" {
		t.Fatalf("missing X-Request-Id header")
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

	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))

	return &App{
        tpls:     tpls,
        staticFS: http.Dir(td),
        cfg: config.Config{
            Name:    "Elliot Alderson",
            Email:   "name@domain.com",
            Brand:   "mr.robot",
            Tagline: "Computer Repair with a Smile",
        },
        log: logger, 
    }

}

func ioReadAll(r io.Reader) (string, error) {
	var sb strings.Builder
	_, err := io.Copy(&sb, r)
	return sb.String(), err
}


// cmd/web/main.go
package main

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/brandondunbar/personal-site/internal/config"
)

type App struct {
	tpls     *template.Template
	staticFS http.FileSystem
	cfg      config.Config
}

type TemplateData struct {
	Site config.Config
	Year int
}

func NewApp() (*App, error) {
	tpls, err := template.ParseFiles(
		templatePath("web/templates/base.html.tmpl"),
		templatePath("web/templates/home.html.tmpl"),
	)
	if err != nil {
		return nil, err
	}

	// Load JSON config
	cfg, err := config.LoadConfig(templatePath("configs/site.json"))
	if err != nil {
		return nil, err
	}

	return &App{
		tpls:     tpls,
		staticFS: http.Dir(templatePath("web/static")),
		cfg:      cfg,
	}, nil
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Static assets with long cache
	fs := http.FileServer(a.staticFS)
	mux.Handle("/static/", cacheControl(http.StripPrefix("/static/", fs)))

	// Home → render template using wrapper data and buffered output
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TemplateData{
			Site: a.cfg,
			Year: time.Now().Year(),
		}

		var buf bytes.Buffer
		if err := a.tpls.ExecuteTemplate(&buf, "base", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(w)
	})

	return mux
}

func cacheControl(next http.Handler) http.Handler {
	const cc = "public, max-age=31536000, immutable"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", cc)
		next.ServeHTTP(w, r)
	})
}

func templatePath(rel string) string {
	_, file, _, _ := runtime.Caller(0) // file = .../cmd/web/main.go
	return filepath.Join(filepath.Dir(file), "..", "..", rel)
}

func main() {
	addr := ":8080"
	println("Server listening on http://localhost" + addr)

	app, err := NewApp()
	if err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(addr, app.routes()); err != nil {
		panic(err)
	}
}

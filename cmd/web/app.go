// cmd/web/app.go
package main

import (
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

// Small helpers kept with the App since they're used by wiring.

func cacheControl(next http.Handler) http.Handler {
	const cc = "public, max-age=31536000, immutable"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", cc)
		next.ServeHTTP(w, r)
	})
}

func templatePath(rel string) string {
	_, file, _, _ := runtime.Caller(0) // file = .../cmd/web/app.go (this file)
	return filepath.Join(filepath.Dir(file), "..", "..", rel)
}

// now() kept as a function for easy stubbing later if needed.
func now() time.Time { return time.Now() }


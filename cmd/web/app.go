
// cmd/web/app.go
package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/brandondunbar/personal-site/internal/config"
)

type App struct {
	tpls     *template.Template
	staticFS http.FileSystem
	cfg      config.Config
	log      *slog.Logger
}

type TemplateData struct {
	Site config.Config
	Year int
}

func NewApp() (*App, error) {
	tpls, err := template.ParseFiles(
		templatePath("web/templates/base.html.tmpl"),
		templatePath("web/templates/home.html.tmpl"),
		// Optional if you have them:
		templatePath("web/templates/404.html.tmpl"),
		templatePath("web/templates/500.html.tmpl"),
	)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadConfig(templatePath("configs/site.json"))
	if err != nil {
		return nil, err
	}

	return &App{
		tpls:     tpls,
		staticFS: http.Dir(templatePath("web/static")),
		cfg:      cfg,
		log:      newLogger(),
	}, nil
}

func newLogger() *slog.Logger {
	var h slog.Handler
	if os.Getenv("LOG_FORMAT") == "json" || os.Getenv("APP_ENV") == "prod" {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.New(h).With(slog.String("service", "personal-site"))
}


func cacheControl(next http.Handler) http.Handler {
	const cc = "public, max-age=31536000, immutable"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", cc)
		next.ServeHTTP(w, r)
	})
}

func templatePath(rel string) string {
	_, file, _, _ := runtime.Caller(0) // this file's path
	return filepath.Join(filepath.Dir(file), "..", "..", rel)
}

func now() time.Time { return time.Now() }


// cmd/web/routes.go
package main

import (
	"bytes"
	"net/http"

	"github.com/brandondunbar/personal-site/internal/httpx"
)

func (a *App) Routes() http.Handler {
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

	// Home — only for "/"; others → custom 404
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			a.renderNotFound(w, r)
			return
		}
		data := TemplateData{Site: a.cfg, Year: now().Year()}
		var buf bytes.Buffer
		if err := a.tpls.ExecuteTemplate(&buf, "base", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(w)
	})

	// Chain: RequestID -> Recover -> Logger
	h := httpx.RequestID(mux)
	h = a.recoverMiddleware(h)
	h = httpx.Logger(a.log)(h)
	return h
}

// custom 404
func (a *App) renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	if a != nil && a.tpls != nil {
		data := TemplateData{Site: a.cfg, Year: now().Year()}
		if err := a.tpls.ExecuteTemplate(w, "notfound", data); err == nil {
			return
		}
	}
	_, _ = w.Write([]byte(`<!doctype html><meta charset="utf-8"><title>Not Found</title><h1>Page not found</h1>`))
}


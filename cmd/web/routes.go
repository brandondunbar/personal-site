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

	// Home — only for "/"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data := TemplateData{
			Site: a.cfg,
			Year: now().Year(),
		}
		var buf bytes.Buffer
		if err := a.tpls.ExecuteTemplate(&buf, "base", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(w)
	})

	// Middleware chain: RequestID -> Recover -> Logger
	h := httpx.RequestID(mux)   // sets X-Request-Id and stores in context
	h = a.recoverMiddleware(h)  // friendly 500 + panic logging
	h = httpx.Logger(a.log)(h)  // JSON access log with request_id, method, path, status, duration

	return h
}


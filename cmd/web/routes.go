// cmd/web/routes.go
package main

import (
	"bytes"
	"net/http"
)

// Routes constructs and returns the HTTP handler tree for the app.
// Keeping this exported makes tests (and future binaries) able to import and verify routing in isolation.
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

	// Home → render template using wrapper data and buffered output
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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

	return mux
}


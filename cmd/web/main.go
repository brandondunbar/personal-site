// cmd/web/main.go
package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
)

// Templates (resolved relative to repo root)
var templates = template.Must(template.ParseFiles(templatePath("web/templates/home.html.tmpl")))

// staticFS is the filesystem used to serve /static.
// Tests can overwrite this to point at a temp dir.
var staticFS http.FileSystem = http.Dir(templatePath("web/static"))

func templatePath(rel string) string {
	// Resolve path relative to this source file's directory.
	_, file, _, _ := runtime.Caller(0) // file = .../cmd/web/main.go
	return filepath.Join(filepath.Dir(file), "..", "..", rel)
}

func routes() http.Handler {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Static assets with long cache
	// FileServer sets Last-Modified; we add Cache-Control.
	fs := http.FileServer(staticFS)
	mux.Handle("/static/", cacheControl(http.StripPrefix("/static/", fs)))

	// Home → render template
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := templates.ExecuteTemplate(w, "home.html.tmpl", nil); err != nil {
			http.Error(w, "template error", http.StatusInternalServerError)
		}
	})

	return mux
}

func cacheControl(next http.Handler) http.Handler {
	// Use a long cache; prefer filename fingerprinting for prod to make it safe to be immutable.
	const cc = "public, max-age=31536000, immutable"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", cc)
		next.ServeHTTP(w, r)
	})
}

func main() {
	addr := ":8080"
	println("Server listening on http://localhost" + addr)
	if err := http.ListenAndServe(addr, routes()); err != nil {
		panic(err)
	}
}

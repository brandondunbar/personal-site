// cmd/web/middleware.go
package main

import (
	"log"
	"net/http"
	"runtime/debug"
)

func (a *App) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log panic + stack
				log.Printf("panic: %v\n%s", rec, debug.Stack())

				// Always return a friendly 500 page
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)

				// Try rendering a dedicated template; fall back to a small HTML snippet.
				if a != nil && a.tpls != nil {
					data := TemplateData{Site: a.cfg, Year: now().Year()}
					if err := a.tpls.ExecuteTemplate(w, "error500", data); err == nil {
						return
					}
				}
				_, _ = w.Write([]byte(`<!doctype html><meta charset="utf-8">
<title>Something went wrong</title>
<h1>We hit a snag</h1>
<p>Sorry about that—please try again.</p>`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}


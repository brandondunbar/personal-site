// cmd/web/middleware.go
package main

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (a *App) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				if a.log != nil {
					a.log.Error("panic",
						slog.Any("err", rec),
						slog.String("stack", string(debug.Stack())),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					)
				}
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				// Try template; fall back to simple HTML.
				if a != nil && a.tpls != nil {
					data := TemplateData{Site: a.cfg, Year: now().Year()}
					if err := a.tpls.ExecuteTemplate(w, "error500", data); err == nil {
						return
					}
				}
				_, _ = w.Write([]byte(`<!doctype html><meta charset="utf-8"><title>Something went wrong</title><h1>We hit a snag</h1>`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}


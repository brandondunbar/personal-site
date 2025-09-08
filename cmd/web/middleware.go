// cmd/web/middleware.go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

type ctxKey string

const ctxKeyReqID ctxKey = "req_id"

func (a *App) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log panic + stack with request context
				if a.log != nil {
					a.log.Error("panic",
						slog.Any("err", rec),
						slog.String("stack", string(debug.Stack())),
						slog.String("req_id", reqIDFromCtx(r.Context())),
						slog.String("method", r.Method),
						slog.String("path", r.URL.Path),
					)
				}

				// Friendly 500 page
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
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

type statusWriter struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}

func reqIDFromCtx(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyReqID).(string); ok && v != "" {
		return v
	}
	return ""
}

func newReqID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return time.Now().UTC().Format("20060102T150405.000000000Z07:00")
	}
	return hex.EncodeToString(b[:])
}

func remoteIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}


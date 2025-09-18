// cmd/web/routes.go
package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/brandondunbar/personal-site/internal/blog"
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

	mux.HandleFunc("/_test/500", func(w http.ResponseWriter, r *http.Request) {
          appErr := fmt.Errorf("simulated failure for 500 test")
          a.renderServerError(w, r, appErr)
        })
        mux.HandleFunc("/_test/panic", func(w http.ResponseWriter, r *http.Request) {
          panic("boom")
        })

	// Static assets with long cache
	fs := http.FileServer(a.staticFS)
	mux.Handle("/static/", cacheControl(http.StripPrefix("/static/", fs)))

	// Blog index
	mux.HandleFunc("/blog", func(w http.ResponseWriter, r *http.Request) {
		posts := a.blog.All()
		data := struct {
			TemplateData
			Posts []blog.Post
		}{
			TemplateData: TemplateData{Site: a.cfg, Year: now().Year()},
			Posts:        posts,
		}
		var buf bytes.Buffer
		if err := a.tpls.ExecuteTemplate(&buf, "blog_index", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(w)
	})

	// Blog detail /blog/{slug}
	mux.HandleFunc("/blog/", func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/blog/")
		if slug == "" || strings.Contains(slug, "/") {
			a.renderNotFound(w, r)
			return
		}
		post, ok := a.blog.BySlug(slug)
		if !ok {
			a.renderNotFound(w, r)
			return
		}
		data := struct {
			TemplateData
			Post blog.Post
		}{
			TemplateData: TemplateData{Site: a.cfg, Year: now().Year()},
			Post:         post,
		}
		var buf bytes.Buffer
		if err := a.tpls.ExecuteTemplate(&buf, "blog_post", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(w)
	})

	// Home — only for "/"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			a.renderNotFound(w, r)
			return
		}
		data := TemplateData{Site: a.cfg, Year: now().Year(), Title: "Home | " + a.cfg.Title} 
		var buf bytes.Buffer
		if err := a.tpls.ExecuteTemplate(&buf, "home", data); err != nil {
			http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = buf.WriteTo(w)
	})

	// Middleware chain: RequestID -> Recover -> Logger
	h := httpx.RequestID(mux)
	h = a.recoverMiddleware(h)
	h = httpx.Logger(a.log)(h)
	return h
}

// renderNotFound renders a custom 404 page if template "notfound" exists; otherwise a tiny HTML fallback.
func (a *App) renderNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)

	if a != nil && a.tpls != nil {
		data := TemplateData{Site: a.cfg, Year: now().Year(), Title: "Not Found | " + a.cfg.Title}
		if err := a.tpls.ExecuteTemplate(w, "notfound", data); err == nil {
			return
		}
	}
	_, _ = w.Write([]byte(`<!doctype html><meta charset="utf-8"><title>Not Found</title><h1>Page not found</h1>`))
}

func (a *App) renderServerError(w http.ResponseWriter, r *http.Request, err error) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusInternalServerError)

    if a != nil && a.tpls != nil {
        data := TemplateData{
            Site:  a.cfg,
            Year:  now().Year(),
            Title: "Server Error | " + a.cfg.Title,
        }
        if tplErr := a.tpls.ExecuteTemplate(w, "servererror", data); tplErr == nil {
            return
        }
    }

    // fallback plain text
    _, _ = w.Write([]byte(`<!doctype html><meta charset="utf-8"><title>Server Error</title><h1>500 — Internal Server Error</h1>`))
}


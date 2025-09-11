// internal/blog/types.go
package blog

import (
	"html/template"
	"time"
)

type Post struct {
	Title   string
	Slug    string        // url id, e.g. "go-stdlib-web"
	Date    time.Time
	Tags    []string
	Draft   bool
	Summary string
	HTML    template.HTML // rendered markdown
}

type Store interface {
	All() []Post               // sorted desc by Date, no drafts (unless configured)
	BySlug(slug string) (Post, bool)
	ByTag(tag string) []Post   // optional; can stub for now
}


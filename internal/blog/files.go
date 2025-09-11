// internal/blog/files.go
package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

/*
Filesystem-backed implementation of the blog Store.

- Reads Markdown files with YAML front matter from a directory.
- Filters drafts/future-dated posts unless configured to show drafts.
- Renders Markdown to HTML using goldmark.
- Ensures unique slugs and provides fast slug lookup.
*/

type FilesStore struct {
	dir        string
	showDrafts bool
	now        func() time.Time

	posts  []Post
	bySlug map[string]int
}

// Functional options
type FilesOption func(*FilesStore)

// WithDrafts enables returning draft/future posts (useful in dev).
func WithDrafts(v bool) FilesOption { return func(s *FilesStore) { s.showDrafts = v } }

// WithNow overrides the time source (useful for tests).
func WithNow(f func() time.Time) FilesOption { return func(s *FilesStore) { s.now = f } }

// NewFilesStore loads posts from dir and prepares indexes.
func NewFilesStore(dir string, opts ...FilesOption) (*FilesStore, error) {
	s := &FilesStore{
		dir: dir,
		now: time.Now,
	}
	for _, opt := range opts {
		opt(s)
	}
	if err := s.reload(); err != nil {
		return nil, err
	}
	return s, nil
}

// All returns all posts sorted by date desc (copy).
func (s *FilesStore) All() []Post {
	out := make([]Post, len(s.posts))
	copy(out, s.posts)
	return out
}

// BySlug returns a post by its slug.
func (s *FilesStore) BySlug(slug string) (Post, bool) {
	i, ok := s.bySlug[slug]
	if !ok {
		return Post{}, false
	}
	return s.posts[i], true
}

// ByTag returns posts with a given tag (case-insensitive), sorted by date desc.
func (s *FilesStore) ByTag(tag string) []Post {
	tag = strings.ToLower(tag)
	var out []Post
	for _, p := range s.posts {
		for _, t := range p.Tags {
			if strings.ToLower(t) == tag {
				out = append(out, p)
				break
			}
		}
	}
	return out
}

/************ loading ************/

func (s *FilesStore) reload() error {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return err
	}

	var posts []Post
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		p, ok, err := parseFile(filepath.Join(s.dir, e.Name()), s.showDrafts, s.now)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		posts = append(posts, p)
	}

	// Sort newest first; zero dates go last.
	sort.SliceStable(posts, func(i, j int) bool {
		di, dj := posts[i].Date, posts[j].Date
		if di.IsZero() && dj.IsZero() {
			return posts[i].Title < posts[j].Title
		}
		if di.IsZero() {
			return false
		}
		if dj.IsZero() {
			return true
		}
		return di.After(dj)
	})

	// Ensure slug uniqueness by appending -2, -3, ...
	seen := make(map[string]struct{}, len(posts))
	for i := range posts {
		base := posts[i].Slug
		slug := base
		k := 2
		for {
			if _, exists := seen[slug]; !exists {
				break
			}
			slug = fmt.Sprintf("%s-%d", base, k)
			k++
		}
		posts[i].Slug = slug
		seen[slug] = struct{}{}
	}

	s.posts = posts
	s.bySlug = make(map[string]int, len(posts))
	for i, p := range posts {
		s.bySlug[p.Slug] = i
	}
	return nil
}

/************ parsing ************/

var md = goldmark.New() // customize later with extensions if needed

type frontMatter struct {
	Title   string   `yaml:"title"`
	Slug    string   `yaml:"slug"`
	Date    string   `yaml:"date"`
	Tags    []string `yaml:"tags"`
	Draft   bool     `yaml:"draft"`
	Summary string   `yaml:"summary"`
}

// parseFile reads a .md file, parses front matter, renders Markdown.
// Returns (Post, true, nil) when included;
// (zero, false, nil) when excluded due to draft/future;
// (zero, false, err) on error.
func parseFile(path string, showDrafts bool, now func() time.Time) (Post, bool, error) {
	var zero Post

	b, err := os.ReadFile(path)
	if err != nil {
		return zero, false, fmt.Errorf("read %s: %w", filepath.Base(path), err)
	}

	fmBytes, body := splitFrontMatter(b)
	var fm frontMatter
	if len(fmBytes) > 0 {
		if err := yaml.Unmarshal(fmBytes, &fm); err != nil {
			return zero, false, fmt.Errorf("front matter %s: %w", filepath.Base(path), err)
		}
	}

	title := strings.TrimSpace(fm.Title)
	if title == "" {
		title = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	}

	slug := fm.Slug
	if slug == "" {
		slug = Slugify(title)
	}

	date, _ := parseDate(fm.Date)

	// Filter drafts/future posts unless showing drafts.
	if !showDrafts {
		if fm.Draft {
			return zero, false, nil
		}
		if !date.IsZero() && date.After(now()) {
			return zero, false, nil
		}
	}

	var out bytes.Buffer
	if err := md.Convert(body, &out); err != nil {
		return zero, false, fmt.Errorf("markdown %s: %w", filepath.Base(path), err)
	}

	post := Post{
		Title:   title,
		Slug:    slug,
		Date:    date,
		Tags:    fm.Tags,
		Draft:   fm.Draft,
		Summary: fm.Summary,
		HTML:    template.HTML(out.String()),
	}
	return post, true, nil
}

// splitFrontMatter returns (yaml, body) if file begins with a '---' line; otherwise (nil, b).
// Supports the common YAML front matter block delimited by:
// ---\n
// ...yaml...
// ---\n
func splitFrontMatter(b []byte) (yamlPart, body []byte) {
	s := string(b)
	// Normalize line endings lightly by working line-by-line.
	lines := strings.Split(s, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return nil, b
	}
	// Find closing '---' on its own line.
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end == -1 {
		return nil, b
	}
	yamlPart = []byte(strings.Join(lines[1:end], "\n"))
	body = []byte(strings.Join(lines[end+1:], "\n"))
	return yamlPart, body
}

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, nil
	}
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02T15:04",
	}
	var last error
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		} else {
			last = err
		}
	}
	return time.Time{}, last
}

// Slugify converts a string into a URL-safe identifier.
// - Lowercases
// - Non-alphanumerics -> single hyphen
// - Trims leading/trailing hyphens
func Slugify(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	prevHyphen := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevHyphen = false
			continue
		}
		if !prevHyphen {
			b.WriteByte('-')
			prevHyphen = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		out = "post"
	}
	return out
}


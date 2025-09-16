package config

import (
	"encoding/json"
	"net"
	"os"
	"strings"
	"time"
)

/* ---------- Site content (backward-compatible API) ---------- */

type Preload struct {
	Href string `json:"Href"`
	As   string `json:"As"`
}

type Head struct {
	MetaDescription string    `json:"MetaDescription"`
	Preloads        []Preload `json:"Preloads"`
	Styles          []string  `json:"Styles"`
	Scripts         []string  `json:"Scripts"`
}

type NavItem struct {
	Href  string `json:"Href"`
	Label string `json:"Label"`
	Class string `json:"Class,omitempty"`
}

type Header struct {
	Brand string    `json:"Brand"`
	Nav   []NavItem `json:"Nav"`
}

type Action struct {
	Href  string `json:"Href"`
	Label string `json:"Label"`
	Class string `json:"Class,omitempty"`
}

type Hero struct {
	Title   string   `json:"Title"`
	Lede    string   `json:"Lede"`
	Actions []Action `json:"Actions"`
}

/* ---------- NEW Work schema (preferred) ---------- */

type ProjectLink struct {
	Href     string `json:"Href"`
	Label    string `json:"Label"`
	Disabled bool   `json:"Disabled,omitempty"`
}

type Project struct {
	Title     string        `json:"Title"`
	Blurb     string        `json:"Blurb,omitempty"` // Preferred short description
	Body      string        `json:"Body,omitempty"`  // Legacy fallback -> mapped into Blurb if Blurb empty
	Role      string        `json:"Role,omitempty"`
	Year      int           `json:"Year,omitempty"`
	Tech      []string      `json:"Tech,omitempty"`      // Preferred: array of tech badges
	Links     []ProjectLink `json:"Links,omitempty"`     // Preferred: array of links
	Thumb     string        `json:"Thumb,omitempty"`
	Images    []string      `json:"Images,omitempty"`
	Highlights []string     `json:"Highlights,omitempty"`
}

/* ---------- Legacy Work schema (auto-mapped) ---------- */

type WorkCard struct {
	Title    string `json:"Title"`
	Tech     string `json:"Tech"` // Legacy: single string ("Go · Postgres · Docker")
	Body     string `json:"Body"`
	LinkHref string `json:"LinkHref"`
	LinkText string `json:"LinkText"`
	Disabled bool   `json:"Disabled,omitempty"`
}

/* ---------- Work container supporting both schemas ---------- */

type Work struct {
	Title    string    `json:"Title"`
	Intro    string    `json:"Intro,omitempty"`
	Projects []Project `json:"Projects,omitempty"` // NEW schema
	// Legacy field still accepted in JSON; mapped in UnmarshalJSON:
	Cards []WorkCard `json:"Cards,omitempty"`
}

func (w *Work) UnmarshalJSON(b []byte) error {
	// Shadow type to avoid recursion
	type rawWork Work
	var rw rawWork
	if err := json.Unmarshal(b, &rw); err != nil {
		return err
	}

	// If Projects already provided (new schema), prefer it
	if len(rw.Projects) > 0 {
		*w = Work(rw)
		return nil
	}

	// Otherwise, map legacy Cards -> Projects
	projs := make([]Project, 0, len(rw.Cards))
	for _, c := range rw.Cards {
		p := Project{
			Title: c.Title,
			// Prefer Blurb; legacy uses Body, so map Body -> Blurb
			Blurb: c.Body,
			Tech:  splitTechString(c.Tech),
			Links: []ProjectLink{
				{
					Href:     c.LinkHref,
					Label:    c.LinkText,
					Disabled: c.Disabled,
				},
			},
			// Thumb/Images/Highlights/Role/Year left empty (legacy had none)
		}
		projs = append(projs, p)
	}

	w.Title = rw.Title
	w.Intro = rw.Intro
	w.Projects = projs
	// Keep legacy Cards populated too (harmless), though not required:
	w.Cards = rw.Cards
	return nil
}

func splitTechString(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	// Accept common separators: dot (·), comma, pipe, slash
	seps := []string{"·", ",", "|", "/"}
	// Normalize all seps to comma
	for _, sep := range seps {
		s = strings.ReplaceAll(s, sep, ",")
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

type ServicesColumn struct {
	Title string   `json:"Title"`
	Items []string `json:"Items"`
}

type Services struct {
	Title   string           `json:"Title"`
	Columns []ServicesColumn `json:"Columns"`
	Alt     bool             `json:"Alt,omitempty"`
}

type Step struct {
	Strong string `json:"Strong"`
	Text   string `json:"Text"`
}

type Links struct {
	GitHub   string `json:"GitHub"`
	LinkedIn string `json:"LinkedIn"`
	CV       string `json:"CV"`
	LeetCode string `json:"LeetCode"`
}

type Book struct {
	Title     string `json:"Title"`
	Author    string `json:"Author"`
	Cover     string `json:"Cover"`
	URL       string `json:"URL"`
	IsReading bool   `json:"IsReading"`
}

type Bookshelf struct {
	ID    string `json:"ID"`
	Title string `json:"Title"`
	Books []Book `json:"Books"`
}


type About struct {
	Title      string   `json:"Title"`
	Paragraphs []string `json:"Paragraphs"`
}

type Contact struct {
	Title      string `json:"Title"`
	EmailLabel string `json:"EmailLabel"`
	Button     Action `json:"Button"`
}

type Footer struct {
	Note string `json:"Note"`
}

type Config struct {
	Title   string `json:"Title"`
	Name    string `json:"Name"`
	Email   string `json:"Email"`
	Brand   string `json:"Brand"`
	Tagline string `json:"Tagline"`

	Head     Head     `json:"Head"`
	Header   Header   `json:"Header"`
	Hero     Hero     `json:"Hero"`
	Work     Work     `json:"Work"`
	Services Services `json:"Services"`
	About    About    `json:"About"`
	Links    Links    `json:"Links"`
	Contact  Contact  `json:"Contact"`
	Footer   Footer   `json:"Footer"`
	Bookshelf Bookshelf `json:"Bookshelf"`
}


// LoadConfig reads site content JSON and applies env overrides.
func LoadConfig(path string) (Config, error) {
	var c Config
	b, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	if v := os.Getenv("SITE_EMAIL"); v != "" {
		c.Email = v
	}
	return c, nil
}

/* ---------- Runtime (env) configuration ---------- */

// Runtime holds process/runtime settings derived from environment.
type Runtime struct {
	Env     string // "dev" or "prod"
	Addr    string // listen address, e.g. ":8080" or "127.0.0.1:9090"
	BaseURL string // externally visible base URL (no trailing slash)
}

// LoadRuntime loads runtime config from env with sane defaults.
// Precedence:
//   ADDR > PORT > default(":8080")
//   BASE_URL used if set, otherwise derived from Addr.
//   Env from APP_ENV|GO_ENV|ENV, normalized to "dev" or "prod".
func LoadRuntime() Runtime {
	env := normalizeEnv(firstNonEmpty(
		os.Getenv("APP_ENV"),
		os.Getenv("GO_ENV"),
		os.Getenv("ENV"),
	))

	addr := os.Getenv("ADDR")
	if addr == "" {
		if p := os.Getenv("PORT"); p != "" {
			addr = ":" + strings.TrimPrefix(p, ":")
		} else {
			addr = ":8080"
		}
	}

	base := strings.TrimRight(os.Getenv("BASE_URL"), "/")
	if base == "" {
		base = deriveBaseURL(addr, env)
	}

	return Runtime{
		Env:     env,
		Addr:    addr,
		BaseURL: base,
	}
}

/* ---------- helpers ---------- */

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func normalizeEnv(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "prod", "production", "release":
		return "prod"
	default:
		return "dev"
	}
}

func deriveBaseURL(addr, env string) string {
	// Scheme: keep it simple; use HTTP unless user provides BASE_URL.
	scheme := "http"
	host := "localhost"
	port := ""

	if strings.HasPrefix(addr, ":") {
		port = strings.TrimPrefix(addr, ":")
	} else {
		// Try host:port; if it fails, just use the raw addr as host (no port).
		if h, p, err := net.SplitHostPort(addr); err == nil {
			host, port = h, p
		} else {
			host = addr
		}
	}

	// Replace wildcard bind addresses with localhost for a usable URL.
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "localhost"
	}

	if port != "" {
		return scheme + "://" + host + ":" + port
	}
	return scheme + "://" + host
}

// (Keep a tiny dependency on time so Runtime can be extended in future without import churn.)
var _ = time.Now


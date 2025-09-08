package config

import (
	"encoding/json"
	"net"
	"os"
	"strings"
	"time"
)

/* ---------- Site content (unchanged API) ---------- */

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

type WorkCard struct {
	Title    string `json:"Title"`
	Tech     string `json:"Tech"`
	Body     string `json:"Body"`
	LinkHref string `json:"LinkHref"`
	LinkText string `json:"LinkText"`
	Disabled bool   `json:"Disabled,omitempty"`
}

type Work struct {
	Title string     `json:"Title"`
	Cards []WorkCard `json:"Cards"`
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

type Approach struct {
	Title string `json:"Title"`
	Steps []Step `json:"Steps"`
}

type Links struct {
	GitHub   string `json:"GitHub"`
	LinkedIn string `json:"LinkedIn"`
	CV       string `json:"CV"`
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
	Approach Approach `json:"Approach"`
	About    About    `json:"About"`
	Links    Links    `json:"Links"`
	Contact  Contact  `json:"Contact"`
	Footer   Footer   `json:"Footer"`
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


package config

import (
	"encoding/json"
	"os"
)

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
	Alt     bool             `json:"Alt,omitempty"` // for section--alt class toggle if you want it
}

type Step struct {
	Strong string `json:"Strong"` // bold part
	Text   string `json:"Text"`   // rest of sentence
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
	// If you prefer explicit label+href entries instead of Links map, use []Action
}

type Contact struct {
	Title      string `json:"Title"`
	EmailLabel string `json:"EmailLabel"` // e.g., "Email me:"
	Button     Action `json:"Button"`
}

type Footer struct {
	Note string `json:"Note"` // trailing text after © YEAR Name
}

type Config struct {
	// Site-wide
	Title   string `json:"Title"` // e.g., "Brandon — Pragmatic Backend & AI"
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

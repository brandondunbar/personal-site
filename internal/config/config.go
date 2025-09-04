// internal/config/config.go
package config

import (
	"encoding/json"
	"os"
)

type NavItem struct {
	Href  string `json:"Href"`
	Label string `json:"Label"`
	Class string `json:"Class,omitempty"`
}

type Links struct {
	GitHub   string `json:"GitHub"`
	LinkedIn string `json:"LinkedIn"`
	CV       string `json:"CV"`
}

type Config struct {
	Name    string    `json:"Name"`
	Email   string    `json:"Email"`
	Brand   string    `json:"Brand"`
	Tagline string    `json:"Tagline"`
	Links   Links     `json:"Links"`
	Nav     []NavItem `json:"Nav"`
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
	// Optional: environment overrides, e.g. if EMAIL is set
	if v := os.Getenv("SITE_EMAIL"); v != "" {
		c.Email = v
	}
	return c, nil
}

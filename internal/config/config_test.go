package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("LoadRuntime_defaults", func(t *testing.T) {
		// isolate env
		t.Setenv("APP_ENV", "")
		t.Setenv("GO_ENV", "")
		t.Setenv("ENV", "")
		t.Setenv("ADDR", "")
		t.Setenv("PORT", "")
		t.Setenv("BASE_URL", "")

		rt := LoadRuntime()
		if rt.Env != "dev" {
			t.Fatalf("Env = %q, want dev", rt.Env)
		}
		if rt.Addr != ":8080" {
			t.Fatalf("Addr = %q, want :8080", rt.Addr)
		}
		if rt.BaseURL != "http://localhost:8080" {
			t.Fatalf("BaseURL = %q, want http://localhost:8080", rt.BaseURL)
		}
	})

	t.Run("LoadRuntime_port_only", func(t *testing.T) {
		t.Setenv("PORT", "3000")
		t.Setenv("ADDR", "")
		t.Setenv("BASE_URL", "")
		t.Setenv("APP_ENV", "development")

		rt := LoadRuntime()
		if rt.Addr != ":3000" {
			t.Fatalf("Addr = %q, want :3000", rt.Addr)
		}
		if rt.BaseURL != "http://localhost:3000" {
			t.Fatalf("BaseURL = %q, want http://localhost:3000", rt.BaseURL)
		}
		if rt.Env != "dev" {
			t.Fatalf("Env = %q, want dev", rt.Env)
		}
	})

	t.Run("LoadRuntime_addr_wins", func(t *testing.T) {
		t.Setenv("ADDR", "127.0.0.1:9090")
		t.Setenv("PORT", "1234")
		t.Setenv("BASE_URL", "")
		t.Setenv("APP_ENV", "")

		rt := LoadRuntime()
		if rt.Addr != "127.0.0.1:9090" {
			t.Fatalf("Addr = %q, want 127.0.0.1:9090", rt.Addr)
		}
		if rt.BaseURL != "http://127.0.0.1:9090" {
			t.Fatalf("BaseURL = %q, want http://127.0.0.1:9090", rt.BaseURL)
		}
	})

	t.Run("LoadRuntime_prod_with_BASE_URL", func(t *testing.T) {
		t.Setenv("APP_ENV", "production")
		t.Setenv("BASE_URL", "https://brandon.dev")
		t.Setenv("PORT", "8081")
		t.Setenv("ADDR", "")

		rt := LoadRuntime()
		if rt.Env != "prod" {
			t.Fatalf("Env = %q, want prod", rt.Env)
		}
		if rt.Addr != ":8081" {
			t.Fatalf("Addr = %q, want :8081", rt.Addr)
		}
		if rt.BaseURL != "https://brandon.dev" {
			t.Fatalf("BaseURL = %q, want https://brandon.dev", rt.BaseURL)
		}
	})

	t.Run("LoadConfig_overrides_email_from_env", func(t *testing.T) {
		td := t.TempDir()
		path := filepath.Join(td, "site.json")
		json := `{"Title":"T","Name":"N","Email":"orig@example.com"}`
		if err := os.WriteFile(path, []byte(json), 0o644); err != nil {
			t.Fatalf("write site.json: %v", err)
		}
		t.Setenv("SITE_EMAIL", "override@example.com")

		cfg, err := LoadConfig(path)
		if err != nil {
			t.Fatalf("LoadConfig: %v", err)
		}
		if cfg.Email != "override@example.com" {
			t.Fatalf("Email = %q, want override@example.com", cfg.Email)
		}
	})
}


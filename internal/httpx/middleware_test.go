// internal/httpx/middleware_test.go
package httpx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger_WritesExpectedFields(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Build chain: RequestID(inner) -> test handler -> Logger(outer)
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent) // 204
	}))
	h = Logger(logger)(h)

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	// One log line expected
	line := strings.TrimSpace(buf.String())
	if line == "" {
		t.Fatalf("no log output")
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("invalid json log: %v\nline=%s", err, line)
	}

	// msg is "http_request"
	if got := m["msg"]; got != "http_request" {
		t.Fatalf("msg = %v, want http_request", got)
	}

	// Required structured fields
	if v := m["request_id"]; v == nil || v == "" {
		t.Fatalf("missing request_id field")
	}
	if v := m["method"]; v != "GET" {
		t.Fatalf("method = %v, want GET", v)
	}
	if v := m["path"]; v != "/test" {
		t.Fatalf("path = %v, want /test", v)
	}
	// json numbers become float64
	if v, ok := m["status"].(float64); !ok || int(v) != 204 {
		t.Fatalf("status = %v, want 204", m["status"])
	}
	if v := m["duration"]; v == nil || v == "" {
		t.Fatalf("missing duration")
	}

	// Also ensure header got set by RequestID
	if id := rr.Header().Get("X-Request-Id"); id == "" {
		t.Fatalf("X-Request-Id header not set")
	}
}


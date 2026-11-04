package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"skeleton/pkg/config"
	"skeleton/pkg/logger"

	"github.com/labstack/echo/v4"
)

func TestRouterContracts(t *testing.T) {
	appPath := t.TempDir()
	writeAsset(t, appPath, "index.html", "<!doctype html><title>Skeleton</title>")
	writeAsset(t, appPath, "styles.css", "body{}")

	router := NewRouter(&config.Config{
		Server: config.Server{Host: "127.0.0.1", Port: 38600},
		App:    config.App{Path: appPath, File: "index.html"},
	}, nil)
	router.GET("/fail", func(c echo.Context) error {
		return errors.New("private failure detail")
	})

	t.Run("API info", func(t *testing.T) {
		response := request(router, "/api")
		if response.Code != http.StatusOK {
			t.Fatalf("status = %d", response.Code)
		}
		var body apiInfo
		decodeJSON(t, response, &body)
		if body.Name != "Go API Skeleton" || body.Status != "ok" {
			t.Fatalf("body = %#v", body)
		}
	})

	t.Run("readiness without database", func(t *testing.T) {
		response := request(router, "/health/ready")
		if response.Code != http.StatusOK {
			t.Fatalf("status = %d", response.Code)
		}
	})

	t.Run("static index", func(t *testing.T) {
		response := request(router, "/")
		if response.Code != http.StatusOK || response.Body.String() != "<!doctype html><title>Skeleton</title>" {
			t.Fatalf("status = %d body = %q", response.Code, response.Body.String())
		}
	})

	t.Run("JSON not found", func(t *testing.T) {
		response := request(router, "/missing")
		if response.Code != http.StatusNotFound {
			t.Fatalf("status = %d", response.Code)
		}
		var body errorResponse
		decodeJSON(t, response, &body)
		if body.Error.Code != "not_found" || body.Error.RequestID == "" {
			t.Fatalf("body = %#v", body)
		}
	})

	t.Run("internal error hides details", func(t *testing.T) {
		response := request(router, "/fail")
		if response.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d", response.Code)
		}
		var body errorResponse
		decodeJSON(t, response, &body)
		if body.Error.Code != "internal_server_error" || body.Error.Message != "internal server error" {
			t.Fatalf("body = %#v", body)
		}
	})

	t.Run("security headers", func(t *testing.T) {
		response := request(router, "/api")
		if got := response.Header().Get("X-Content-Type-Options"); got != "nosniff" {
			t.Fatalf("X-Content-Type-Options = %q", got)
		}
		if got := response.Header().Get("Content-Security-Policy"); got == "" {
			t.Fatal("Content-Security-Policy is empty")
		}
		if got := response.Header().Get("X-Request-ID"); got == "" {
			t.Fatal("X-Request-ID is empty")
		}
	})
}

func TestRequestLogsIncludeMiddlewareResponses(t *testing.T) {
	logFile, err := os.CreateTemp(t.TempDir(), "requests-*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = logFile.Close() })
	stdout := os.Stdout
	os.Stdout = logFile
	logger.Init()
	os.Stdout = stdout
	router := NewRouter(&config.Config{App: config.App{Path: t.TempDir(), File: "index.html"}}, nil)
	logger.Init()
	logReader, err := os.Open(logFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = logReader.Close() })

	oversized := httptest.NewRequest(http.MethodGet, "/api", strings.NewReader("x"))
	oversized.ContentLength = 3 << 20
	preflight := httptest.NewRequest(http.MethodOptions, "/api", nil)
	preflight.Header.Set("Origin", "https://example.test")
	preflight.Header.Set("Access-Control-Request-Method", http.MethodGet)
	decoder := json.NewDecoder(logReader)
	for _, tt := range []struct {
		name   string
		req    *http.Request
		status int
	}{
		{"oversized request", oversized, http.StatusRequestEntityTooLarge},
		{"CORS preflight", preflight, http.StatusNoContent},
	} {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			router.ServeHTTP(response, tt.req)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d", response.Code, tt.status)
			}
			var record struct {
				Method    string `json:"method"`
				Status    int    `json:"status"`
				RequestID string `json:"request_id"`
			}
			if err := decoder.Decode(&record); err != nil {
				t.Fatalf("request log missing: %v", err)
			}
			if record.Method != tt.req.Method || record.Status != tt.status || record.RequestID == "" || record.RequestID != response.Header().Get("X-Request-ID") {
				t.Fatalf("request log = %+v", record)
			}
		})
	}
}

func request(handler http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, req)
	return response
}

func decodeJSON(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.Unmarshal(response.Body.Bytes(), target); err != nil {
		t.Fatalf("decode JSON: %v; body = %q", err, response.Body.String())
	}
}

func writeAsset(t *testing.T, dir, name, contents string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(contents), 0o600); err != nil {
		t.Fatalf("write asset: %v", err)
	}
}

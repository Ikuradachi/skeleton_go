package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"skeleton/pkg/config"

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

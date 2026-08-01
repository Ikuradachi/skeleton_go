package config

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadUsesRunnableDefaults(t *testing.T) {
	t.Chdir(t.TempDir())
	unsetConfigEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" || cfg.Server.Port != 38600 {
		t.Fatalf("server defaults = %#v", cfg.Server)
	}
	if cfg.DB.Enabled() {
		t.Fatalf("database should be disabled by default: %#v", cfg.DB)
	}
	if cfg.App.Path != "app" || cfg.App.File != "index.html" {
		t.Fatalf("app defaults = %#v", cfg.App)
	}
}

func TestLoadAppliesEnvironmentAfterFile(t *testing.T) {
	unsetConfigEnv(t)
	configPath := writeConfig(t, `
server:
  host: "127.0.0.1"
  port: 4000
db:
  driver: "mysql"
  dsn: "user:pass@tcp(localhost:3306)/app"
app:
  path: "public"
  file: "home.html"
`)
	t.Setenv("SKELETON_CONFIG_FILE", configPath)
	t.Setenv("SKELETON_SERVER_PORT", "5000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != 5000 {
		t.Fatalf("server port = %d, want 5000", cfg.Server.Port)
	}
	if !cfg.DB.Enabled() {
		t.Fatal("database should be enabled")
	}
	if cfg.App.Path != "public" || cfg.App.File != "home.html" {
		t.Fatalf("app config = %#v", cfg.App)
	}
}

func TestLoadRejectsIncompleteDatabaseConfig(t *testing.T) {
	unsetConfigEnv(t)
	configPath := writeConfig(t, `
db:
  driver: "mysql"
`)
	t.Setenv("SKELETON_CONFIG_FILE", configPath)

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "db.driver and db.dsn must be set together") {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestLoadRejectsUnknownKeys(t *testing.T) {
	unsetConfigEnv(t)
	configPath := writeConfig(t, `
server:
  timeout: "5s"
`)
	t.Setenv("SKELETON_CONFIG_FILE", configPath)

	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "decode config") {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestOpenDatabaseSkipsDisabledConfig(t *testing.T) {
	db, err := OpenDatabase(context.Background(), DB{})
	if err != nil {
		t.Fatalf("OpenDatabase() error = %v", err)
	}
	if db != nil {
		t.Fatal("OpenDatabase() returned non-nil database")
	}
}

func TestOpenDatabaseRejectsUnknownDriver(t *testing.T) {
	db, err := OpenDatabase(context.Background(), DB{Driver: "missing", DSN: "dsn"})
	if err == nil || !strings.Contains(err.Error(), "unknown driver") {
		t.Fatalf("OpenDatabase() error = %v", err)
	}
	if db != nil {
		t.Fatal("OpenDatabase() returned non-nil database")
	}
}

func writeConfig(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func unsetConfigEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"SKELETON_CONFIG_FILE",
		"SKELETON_SERVER_HOST",
		"SKELETON_SERVER_PORT",
		"SKELETON_DB_DRIVER",
		"SKELETON_DB_DSN",
		"SKELETON_APP_PATH",
		"SKELETON_APP_FILE",
	} {
		value, ok := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
		if ok {
			t.Cleanup(func() { _ = os.Setenv(key, value) })
		}
	}
}

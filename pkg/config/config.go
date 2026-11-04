package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const envPrefix = "SKELETON"

type Server struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DB struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

func (db DB) Enabled() bool {
	return db.Driver != "" && db.DSN != ""
}

type App struct {
	Path string `mapstructure:"path"`
	File string `mapstructure:"file"`
}

type Config struct {
	Server Server `mapstructure:"server"`
	DB     DB     `mapstructure:"db"`
	App    App    `mapstructure:"app"`
}

// Load reads optional pkg/config/config.yaml values, then applies SKELETON_*
// environment overrides. Defaults are enough to run the API without a database.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("pkg/config")
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(true)
	v.AutomaticEnv()

	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 38600)
	v.SetDefault("db.driver", "")
	v.SetDefault("db.dsn", "")
	v.SetDefault("app.path", "app")
	v.SetDefault("app.file", "index.html")

	if configFile := strings.TrimSpace(os.Getenv(envPrefix + "_CONFIG_FILE")); configFile != "" {
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("read config file %q: %w", configFile, err)
		}
	} else if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.UnmarshalExact(&cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	cfg.Server.Host = strings.TrimSpace(cfg.Server.Host)
	cfg.DB.Driver = strings.TrimSpace(cfg.DB.Driver)
	cfg.App.Path = strings.TrimSpace(cfg.App.Path)
	cfg.App.File = strings.TrimSpace(cfg.App.File)

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return &cfg, nil
}

func (cfg Config) Validate() error {
	if cfg.Server.Host == "" {
		return errors.New("server.host must not be empty")
	}
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535, got %d", cfg.Server.Port)
	}
	if (cfg.DB.Driver == "") != (cfg.DB.DSN == "") {
		return errors.New("db.driver and db.dsn must be set together")
	}
	if cfg.App.Path == "" {
		return errors.New("app.path must not be empty")
	}
	if cfg.App.File == "" {
		return errors.New("app.file must not be empty")
	}
	return nil
}

// OpenDatabase opens and verifies the configured database. It returns nil when
// database settings are omitted.
func OpenDatabase(ctx context.Context, cfg DB) (*sqlx.DB, error) {
	if !cfg.Enabled() {
		return nil, nil
	}

	db, err := sqlx.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(10)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

package config

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type (
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	DB struct {
		Driver string `mapstructure:"driver"`
		DSN    string `mapstructure:"dsn"`
	}

	App struct {
		Path string `mapstructure:"path"`
		File string `mapstructure:"file"`
	}

	Config struct {
		Server Server `mapstructure:"server"`
		DB     DB     `mapstructure:"db"`
		App    App    `mapstructure:"app"`
	}
)

type (
	ServiceConfig interface {
		Load() (*Config, error)
		Init(cfg *Config) (*sqlx.DB, error)
	}

	serviceConfigImpl struct {
		//
	}
)

func (s *serviceConfigImpl) Load() (*Config, error) {
	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("pkg/config")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (s *serviceConfigImpl) Init(cfg *Config) (*sqlx.DB, error) {
	// Connect to database
	db, err := sqlx.Connect(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func NewServiceConfig() ServiceConfig {
	return &serviceConfigImpl{}
}

package main

import (
	"fmt"
	"skeleton/internal/api"
	"skeleton/pkg/config"

	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.NewServiceConfig()

	// Load configuration
	conf, err := cfg.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		panic(err)
	}

	// Initialize database connection
	db, err := cfg.Init(conf)
	if err != nil {
		fmt.Println("failed to init database:", err)
		panic(err)
	}

	// Initialize API router
	router := api.NewRouter(conf)
	router.Logger.Fatal(
		router.Start(fmt.Sprintf(":%d", conf.Server.Port)),
	)

	// Close database connection when the application exits
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("failed to close database:", err)
		}
	}(db)
}

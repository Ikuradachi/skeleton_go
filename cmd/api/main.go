package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"skeleton/internal/api"
	"skeleton/pkg/config"
	"skeleton/pkg/logger"
	"time"

	"github.com/jmoiron/sqlx"
)

func main() {
	// Initialize logger
	logger.Init()

	cfg := config.NewServiceConfig()

	// Load configuration
	conf, err := cfg.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		panic(err)
	}

	// Initialize database connection
	db, err := cfg.Init(conf)
	if err != nil {
		logger.Error("failed to init database", "error", err)
		panic(err)
	}

	// Close database connection when the application exits
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("failed to close database", "error", err)
		}
	}(db)

	// Initialize API router
	router := api.NewRouter(conf, db)

	// Start server in a goroutine
	go func() {
		logger.Info("starting server", "port", conf.Server.Port)
		if err := router.Start(fmt.Sprintf(":%d", conf.Server.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("shutting down the server", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// Gracefully shutdown the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down server...")
	if err := router.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exited")
}

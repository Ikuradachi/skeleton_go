package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"skeleton/internal/api"
	"skeleton/pkg/config"
	"skeleton/pkg/logger"
)

const (
	databaseConnectTimeout = 5 * time.Second
	shutdownTimeout        = 10 * time.Second
)

func main() {
	logger.Init()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		logger.Error("application stopped", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	dbCtx, cancelDB := context.WithTimeout(ctx, databaseConnectTimeout)
	db, err := config.OpenDatabase(dbCtx, conf.DB)
	cancelDB()
	if err != nil {
		return err
	}
	if db != nil {
		defer func() {
			if err := db.Close(); err != nil {
				logger.Error("close database", "error", err)
			}
		}()
	}

	router := api.NewRouter(conf, db)
	router.Server.ReadHeaderTimeout = 5 * time.Second
	router.Server.ReadTimeout = 15 * time.Second
	router.Server.WriteTimeout = 30 * time.Second
	router.Server.IdleTimeout = 60 * time.Second

	address := net.JoinHostPort(conf.Server.Host, strconv.Itoa(conf.Server.Port))
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("listen HTTP: %w", err)
	}
	defer func() { _ = listener.Close() }()
	router.Listener = listener

	serverErr := make(chan error, 1)
	logger.Info("server listening", "address", listener.Addr().String(), "database", conf.DB.Enabled())
	go func() {
		serverErr <- router.Start(address)
	}()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	case <-ctx.Done():
	}

	logger.Info("server shutting down")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()
	if err := router.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shut down HTTP server: %w", err)
	}

	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve HTTP: %w", err)
	}

	logger.Info("server stopped")
	return nil
}

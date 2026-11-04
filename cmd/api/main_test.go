package main

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRunStopsAfterCancellation(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	address := listener.Addr().String()
	configureRun(t, listener.Addr().(*net.TCPAddr).Port)
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- run(ctx) }()
	client := &http.Client{Timeout: 100 * time.Millisecond, Transport: &http.Transport{Proxy: nil}}
	defer client.CloseIdleConnections()
	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()
	poll := time.NewTicker(10 * time.Millisecond)
	defer poll.Stop()
	for {
		response, err := client.Get("http://" + address + "/health/ready")
		if err == nil {
			_ = response.Body.Close()
			if response.StatusCode != http.StatusOK {
				t.Fatalf("readiness status = %d", response.StatusCode)
			}
			break
		}
		select {
		case err := <-done:
			t.Fatalf("server stopped before ready: %v", err)
		case <-timeout.C:
			t.Fatal("server did not become ready")
		case <-poll.C:
		}
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("shutdown error = %v", err)
		}
	case <-timeout.C:
		t.Fatal("server did not stop after cancellation")
	}
	listener, err = net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("listener not released: %v", err)
	}
	_ = listener.Close()
}

func TestRunRejectsOccupiedPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	configureRun(t, listener.Addr().(*net.TCPAddr).Port)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := run(ctx); err == nil || !strings.Contains(err.Error(), "listen HTTP") {
		t.Fatalf("startup error = %v", err)
	}
}

func TestRunRejectsDatabaseStartupFailure(t *testing.T) {
	configureRun(t, 38600)
	t.Setenv("SKELETON_DB_DRIVER", "missing")
	t.Setenv("SKELETON_DB_DSN", "test")
	if err := run(context.Background()); err == nil || !strings.Contains(err.Error(), "open database") {
		t.Fatalf("startup error = %v", err)
	}
}

func configureRun(t *testing.T, port int) {
	t.Helper()
	t.Chdir(t.TempDir())
	for key, value := range map[string]string{
		"SKELETON_CONFIG_FILE": "",
		"SKELETON_SERVER_HOST": "127.0.0.1",
		"SKELETON_SERVER_PORT": strconv.Itoa(port),
		"SKELETON_DB_DRIVER":   "",
		"SKELETON_DB_DSN":      "",
		"SKELETON_APP_PATH":    t.TempDir(),
		"SKELETON_APP_FILE":    "index.html",
	} {
		t.Setenv(key, value)
	}
}

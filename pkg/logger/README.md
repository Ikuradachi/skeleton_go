# logger

`logger` wraps Go's `log/slog` with a JSON handler writing to stdout at info level.

Initialize it once in the process entrypoint:

```go
logger.Init()
```

Use package helpers for ordinary structured logs:

```go
logger.Info("server listening", "address", ":38600")
logger.Error("database unavailable", "error", err)
```

Use the logger directly when a request context or typed attributes matter:

```go
logger.Get().LogAttrs(ctx, slog.LevelInfo, "request completed",
	slog.String("request_id", requestID),
	slog.Int("status", status),
)
```

Create an independent logger at another level with `NewLogger`:

```go
debugLog := logger.NewLogger(slog.LevelDebug)
```

# Logger Package

A centralized, structured logging utility using Go's standard `log/slog` package.

## Features

- 🎯 **Centralized Logging**: Single logger instance used across the application
- 📊 **Structured Logging**: JSON-formatted logs with key-value pairs
- 🎨 **Multiple Log Levels**: Info, Error, Debug, and Warn
- ⚙️ **Customizable**: Create custom logger instances with different configurations

## Usage

### Initialize the Logger

In your `main.go`:

```go
import "skeleton/pkg/logger"

func main() {
    // Initialize the default logger
    logger.Init()
    
    // Your application code...
}
```

### Using the Logger

#### Option 1: Using Helper Functions (Recommended for Simple Logging)

```go
import "skeleton/pkg/logger"

// Log info messages
logger.Info("server started", "port", 8080)

// Log errors
logger.Error("database connection failed", "error", err)

// Log debug messages
logger.Debug("processing request", "user_id", userId)

// Log warnings
logger.Warn("deprecated API used", "endpoint", "/old-api")
```

#### Option 2: Using the Logger Instance (Recommended for Advanced Usage)

```go
import (
    "log/slog"
    "skeleton/pkg/logger"
)

func MyFunction() {
    log := logger.Get()
    
    log.Info("processing started")
    log.LogAttrs(ctx, slog.LevelInfo, "user action",
        slog.String("user_id", "123"),
        slog.String("action", "login"),
        slog.Int("status", 200),
    )
}
```

### Custom Logger Instance

Create a custom logger with specific log level:

```go
import (
    "log/slog"
    "skeleton/pkg/logger"
)

func main() {
    // Create a debug-level logger
    debugLogger := logger.NewLogger(slog.LevelDebug)
    debugLogger.Debug("this will be logged")
}
```

## Log Output Format

Logs are output in JSON format:

```json
{"time":"2026-02-11T10:30:00Z","level":"INFO","msg":"server started","port":8080}
{"time":"2026-02-11T10:30:01Z","level":"ERROR","msg":"database connection failed","error":"connection timeout"}
```

## Integration with Echo Middleware

The logger is integrated with Echo's RequestLogger middleware for automatic HTTP request logging:

```go
// In internal/api/middleware.go
log := logger.Get()
e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
    // ... configuration
    LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
        log.LogAttrs(c.Request().Context(), slog.LevelInfo, "REQUEST",
            slog.String("uri", v.URI),
            slog.Int("status", v.Status),
        )
        return nil
    },
}))
```

## Best Practices

1. **Initialize Early**: Call `logger.Init()` at the start of your `main()` function
2. **Use Structured Logging**: Always use key-value pairs instead of string concatenation
3. **Choose Appropriate Levels**: 
   - `Info`: General application flow
   - `Error`: Error conditions that need attention
   - `Warn`: Warning conditions
   - `Debug`: Detailed diagnostic information (development only)
4. **Context-Aware**: Use `LogAttrs` with context when available for better tracing

## Future Enhancements

Potential improvements for this logger package:

- [ ] Add configuration from YAML file (log level, format, output destination)
- [ ] Support for multiple output destinations (file, stdout, external services)
- [ ] Add request ID tracking for distributed tracing
- [ ] Performance metrics logging
- [ ] Log rotation support


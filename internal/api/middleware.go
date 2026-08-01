package api

import (
	"log/slog"
	"net/http"

	"skeleton/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewMiddleware(e *echo.Echo) {
	e.Use(middleware.RequestID())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XFrameOptions:      "DENY",
		ContentTypeNosniff: "nosniff",
		XSSProtection:      "0",
		ContentSecurityPolicy: "default-src 'self'; base-uri 'self'; form-action 'self'; " +
			"frame-ancestors 'none'; img-src 'self' data:; object-src 'none'",
		ReferrerPolicy: "strict-origin-when-cross-origin",
	}))
	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{http.MethodGet, http.MethodHead, http.MethodOptions},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
		ExposeHeaders: []string{echo.HeaderXRequestID},
		MaxAge:        300,
	}))

	log := logger.Get()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:       true,
		LogURIPath:      true,
		LogStatus:       true,
		LogLatency:      true,
		LogRemoteIP:     true,
		LogRequestID:    true,
		LogResponseSize: true,
		LogError:        true,
		HandleError:     true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			level := slog.LevelInfo
			message := "request completed"
			if v.Status >= http.StatusInternalServerError {
				level = slog.LevelError
				message = "request failed"
			} else if v.Error != nil {
				level = slog.LevelWarn
				message = "request rejected"
			}
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("path", v.URIPath),
				slog.Int("status", v.Status),
				slog.Duration("latency", v.Latency),
				slog.String("remote_ip", v.RemoteIP),
				slog.String("request_id", v.RequestID),
				slog.Int64("response_size", v.ResponseSize),
			}
			if v.Error != nil {
				attrs = append(attrs, slog.String("error", v.Error.Error()))
			}
			log.LogAttrs(c.Request().Context(), level, message, attrs...)
			return nil
		},
	}))
	e.Use(middleware.Recover())
}

package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewMiddleware(e *echo.Echo) {
	// Enable CORS
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	// Enable request logging
	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))

	// Enable panic recovery
	e.Use(middleware.RecoverWithConfig(middleware.DefaultRecoverConfig))
}

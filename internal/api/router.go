package api

import (
	"net/http"
	"skeleton/pkg/config"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func NewRouter(conf *config.Config, db *sqlx.DB) *echo.Echo {
	e := echo.New()

	// Initialize middleware
	NewMiddleware(e)

	// Initialize static file server
	e.Static("/", conf.App.Path)

	// Initialize API routes
	api := e.Group("/api")
	api.GET("", func(c echo.Context) error {
		return c.String(http.StatusOK, "Api is up and running 🚀")
	})

	_ = db // to avoid unused variable error

	return e
}

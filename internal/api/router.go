package api

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"skeleton/internal/features/health"
	"skeleton/pkg/config"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type apiInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Docs   string `json:"docs"`
}

type apiError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

func NewRouter(conf *config.Config, db *sqlx.DB) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = httpErrorHandler

	NewMiddleware(e)

	var pinger health.Pinger
	if db != nil {
		pinger = db
	}
	healthHandler := health.NewHandler(pinger)
	e.GET("/health/live", healthHandler.Live)
	e.GET("/health/ready", healthHandler.Ready)

	api := e.Group("/api")
	api.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, apiInfo{
			Name:   "Go API Skeleton",
			Status: "ok",
			Docs:   "/#getting-started",
		})
	})

	staticRoot := filepath.Clean(conf.App.Path)
	e.File("/", filepath.Join(staticRoot, conf.App.File))
	e.Static("/", staticRoot)

	return e
}

func httpErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	status := http.StatusInternalServerError
	message := "internal server error"
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		status = httpErr.Code
		if status < http.StatusInternalServerError {
			switch value := httpErr.Message.(type) {
			case string:
				message = value
			case error:
				message = value.Error()
			default:
				message = http.StatusText(status)
			}
		}
	}
	if status == http.StatusNotFound {
		message = "resource not found"
	}

	code := strings.ToLower(strings.ReplaceAll(http.StatusText(status), " ", "_"))
	_ = c.JSON(status, errorResponse{Error: apiError{
		Code:      code,
		Message:   message,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
	}})
}

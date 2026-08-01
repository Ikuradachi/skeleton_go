package health

import (
	"context"
	"net/http"
	"time"

	"skeleton/pkg/logger"

	"github.com/labstack/echo/v4"
)

const readinessTimeout = 2 * time.Second

type Pinger interface {
	PingContext(context.Context) error
}

type Handler struct {
	db Pinger
}

type Response struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks,omitempty"`
}

func NewHandler(db Pinger) *Handler {
	return &Handler{db: db}
}

func (h *Handler) Live(c echo.Context) error {
	return c.JSON(http.StatusOK, Response{Status: "ok"})
}

func (h *Handler) Ready(c echo.Context) error {
	if h.db == nil {
		return c.JSON(http.StatusOK, Response{
			Status: "ok",
			Checks: map[string]string{"database": "disabled"},
		})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), readinessTimeout)
	defer cancel()
	if err := h.db.PingContext(ctx); err != nil {
		logger.Get().WarnContext(c.Request().Context(), "database readiness check failed", "error", err)
		return c.JSON(http.StatusServiceUnavailable, Response{
			Status: "unavailable",
			Checks: map[string]string{"database": "unavailable"},
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status: "ok",
		Checks: map[string]string{"database": "ok"},
	})
}

package health

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type stubPinger struct {
	err error
}

func (p stubPinger) PingContext(context.Context) error {
	return p.err
}

func TestLive(t *testing.T) {
	response := call(t, NewHandler(nil).Live)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	assertHealthResponse(t, response, "ok", nil)
}

func TestReady(t *testing.T) {
	tests := []struct {
		name       string
		pinger     Pinger
		wantStatus int
		wantState  string
		wantDB     string
	}{
		{name: "database disabled", wantStatus: http.StatusOK, wantState: "ok", wantDB: "disabled"},
		{name: "database available", pinger: stubPinger{}, wantStatus: http.StatusOK, wantState: "ok", wantDB: "ok"},
		{name: "database unavailable", pinger: stubPinger{err: errors.New("offline")}, wantStatus: http.StatusServiceUnavailable, wantState: "unavailable", wantDB: "unavailable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := call(t, NewHandler(tt.pinger).Ready)
			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			assertHealthResponse(t, response, tt.wantState, map[string]string{"database": tt.wantDB})
		})
	}
}

func call(t *testing.T, handler echo.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	if err := handler(e.NewContext(req, response)); err != nil {
		t.Fatalf("handler error = %v", err)
	}
	return response
}

func assertHealthResponse(t *testing.T, recorder *httptest.ResponseRecorder, status string, checks map[string]string) {
	t.Helper()
	var response Response
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != status {
		t.Fatalf("response status = %q, want %q", response.Status, status)
	}
	for key, want := range checks {
		if got := response.Checks[key]; got != want {
			t.Fatalf("check %q = %q, want %q", key, got, want)
		}
	}
}

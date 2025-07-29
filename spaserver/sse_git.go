package spaserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) Sse(c echo.Context) error {
	if s.sseManager == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "SSE manager not initialized")
	}

	go func() {
		// Received Browser Disconnection
		<-c.Request().Context().Done()
		s.Logger().Info("The client is disconnected!!!")
	}()
	s.Logger().Info("The client is connected!!!")
	s.sseManager.ServeHTTP(c.Response(), c.Request())

	// Check if response was written successfully
	if c.Response().Committed {
		return nil
	}
	return echo.NewHTTPError(http.StatusInternalServerError, "Failed to establish SSE connection")
}

package handler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v5"
)

// Health repond a un simple contrôle de sante de l'API.
func Health(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "graefik-api",
		"go":      runtime.Version(),
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

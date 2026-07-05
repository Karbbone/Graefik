package server

import (
	"os"
	"strings"

	"github.com/Karbbone/Graefik/apps/api/internal/handler"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// New construit et configure l'instance Echo (middlewares + routes).
func New() *echo.Echo {
	e := echo.New()

	// Middlewares globaux.
	e.Use(middleware.RequestLogger()) // logs structures (slog) avec latence/status
	e.Use(middleware.Recover())       // recupere les panics en erreurs propres
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins(),
	}))

	// Routes de l'API, prefixees par /api.
	api := e.Group("/api")
	api.GET("/health", handler.Health)

	return e
}

// allowedOrigins lit CORS_ORIGINS (liste separee par des virgules) ou autorise tout par defaut.
func allowedOrigins() []string {
	raw := os.Getenv("CORS_ORIGINS")
	if raw == "" {
		return []string{"*"}
	}
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{"*"}
	}
	return origins
}

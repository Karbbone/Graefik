// Package http est un adaptateur INBOUND (pilotant) : il expose le cœur via
// une API HTTP (Echo v5). Il ne connaît que les ports inbound, jamais les
// implémentations concrètes.
package http

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// NewRouter construit l'instance Echo, branche les middlewares et les routes
// sur les services (ports inbound) fournis par la composition root.
func NewRouter(corsOrigins []string, cookie CookieConfig, auth port.AuthService, tasks port.TaskService) *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
	}))

	api := e.Group("/api")

	// Public.
	api.GET("/health", Health)

	authHandler := NewAuthHandler(auth, cookie)
	authGroup := api.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/logout", authHandler.Logout)
	authGroup.GET("/me", authHandler.Me)

	// Protégé par session.
	protected := api.Group("", RequireAuth(auth, cookie.Name))

	taskHandler := NewTaskHandler(tasks)
	protected.GET("/tasks", taskHandler.List)
	protected.POST("/tasks", taskHandler.Create)

	return e
}

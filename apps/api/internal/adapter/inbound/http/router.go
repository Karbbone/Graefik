// Package http est un adaptateur INBOUND (pilotant) : il expose le cœur via
// une API HTTP (Echo v5). Il ne connaît que les ports inbound, jamais les
// implémentations concrètes.
package http

import (
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// NewRouter construit l'instance Echo, branche les middlewares et les routes
// sur les services (ports inbound) fournis par la composition root.
func NewRouter(corsOrigins []string, tasks port.TaskService) *echo.Echo {
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: corsOrigins,
	}))

	api := e.Group("/api")
	api.GET("/health", Health)

	tasksHandler := NewTaskHandler(tasks)
	api.GET("/tasks", tasksHandler.List)
	api.POST("/tasks", tasksHandler.Create)

	return e
}

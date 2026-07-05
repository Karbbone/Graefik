// Command api est la composition root : c'est le SEUL endroit qui connaît les
// implémentations concrètes et câble les dépendances entre elles.
package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"

	adapterhttp "github.com/Karbbone/Graefik/apps/api/internal/adapter/inbound/http"
	"github.com/Karbbone/Graefik/apps/api/internal/adapter/outbound/memory"
	"github.com/Karbbone/Graefik/apps/api/internal/core/service"
	"github.com/Karbbone/Graefik/apps/api/internal/platform"
)

func main() {
	cfg := platform.LoadConfig()

	// Adaptateur outbound (piloté).
	taskRepo := memory.NewTaskRepository()

	// Use case (cœur) : injection du repository, de l'horloge et du générateur d'ID.
	taskService := service.NewTaskService(taskRepo, time.Now, uuid.NewString)

	// Adaptateur inbound (pilotant).
	router := adapterhttp.NewRouter(cfg.CORSOrigins, taskService)

	slog.Info("démarrage de l'API Graefik", "port", cfg.Port, "env", cfg.Env)
	if err := router.Start(":" + cfg.Port); err != nil {
		slog.Error("le serveur s'est arrêté", "error", err)
		os.Exit(1)
	}
}

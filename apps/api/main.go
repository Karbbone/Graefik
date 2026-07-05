package main

import (
	"log/slog"
	"os"

	"github.com/Karbbone/Graefik/apps/api/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := server.New()

	slog.Info("demarrage de l'API Graefik", "port", port, "env", os.Getenv("APP_ENV"))
	if err := e.Start(":" + port); err != nil {
		slog.Error("le serveur s'est arrete", "error", err)
		os.Exit(1)
	}
}

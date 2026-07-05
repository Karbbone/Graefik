// Command api est la composition root : c'est le SEUL endroit qui connaît les
// implémentations concrètes et câble les dépendances entre elles.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"

	adapterhttp "github.com/Karbbone/Graefik/apps/api/internal/adapter/inbound/http"
	"github.com/Karbbone/Graefik/apps/api/internal/adapter/outbound/memory"
	"github.com/Karbbone/Graefik/apps/api/internal/adapter/outbound/security"
	"github.com/Karbbone/Graefik/apps/api/internal/adapter/outbound/sqlite"
	"github.com/Karbbone/Graefik/apps/api/internal/core/service"
	"github.com/Karbbone/Graefik/apps/api/internal/platform"
)

func main() {
	cfg := platform.LoadConfig()

	// --- Adaptateurs outbound (pilotés) ---
	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		slog.Error("ouverture de la base SQLite impossible", "path", cfg.DBPath, "error", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepo := sqlite.NewUserRepository(db)
	sessionRepo := sqlite.NewSessionRepository(db)
	hasher := security.NewBcryptHasher()
	taskRepo := memory.NewTaskRepository()

	// --- Use cases (cœur) ---
	authService := service.NewAuthService(userRepo, sessionRepo, hasher, service.AuthConfig{
		Now:           time.Now,
		NewID:         uuid.NewString,
		NewToken:      platform.GenerateToken,
		NewPassword:   platform.GeneratePassword,
		TTL:           cfg.SessionTTL,
		AdminUsername: "graefik",
	})
	taskService := service.NewTaskService(taskRepo, time.Now, uuid.NewString)

	// --- Bootstrap : admin initial (mot de passe affiché une fois dans les logs) ---
	res, err := authService.EnsureAdmin(context.Background())
	if err != nil {
		slog.Error("initialisation de l'admin impossible", "error", err)
		os.Exit(1)
	}
	if res.Created {
		printAdminBanner(res.Username, res.GeneratedPassword)
	}

	// --- Adaptateur inbound (pilotant) ---
	router := adapterhttp.NewRouter(
		cfg.CORSOrigins,
		adapterhttp.CookieConfig{Name: cfg.CookieName, Secure: cfg.CookieSecure, TTL: cfg.SessionTTL},
		authService,
		taskService,
	)

	slog.Info("démarrage de l'API Graefik", "port", cfg.Port, "env", cfg.Env)
	if err := router.Start(":" + cfg.Port); err != nil {
		slog.Error("le serveur s'est arrêté", "error", err)
		os.Exit(1)
	}
}

// printAdminBanner affiche le compte admin initial dans les logs (une seule fois).
func printAdminBanner(username, password string) {
	const line = "============================================================"
	fmt.Fprintln(os.Stdout, "\n"+line)
	fmt.Fprintln(os.Stdout, "  Graefik — compte administrateur initial")
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintf(os.Stdout, "  Identifiant  : %s\n", username)
	fmt.Fprintf(os.Stdout, "  Mot de passe : %s\n", password)
	fmt.Fprintln(os.Stdout, "")
	fmt.Fprintln(os.Stdout, "  Ce mot de passe n'est affiché qu'UNE seule fois.")
	fmt.Fprintln(os.Stdout, "  Notez-le et changez-le dès que possible.")
	fmt.Fprintf(os.Stdout, "%s\n\n", line)
}

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
	if err := run(); err != nil {
		slog.Error("l'API s'est arrêtée", "error", err)
		os.Exit(1)
	}
}

// run câble les dépendances et démarre le serveur. Isolé de main() pour que les
// defer (fermeture de la BDD) s'exécutent avant tout os.Exit.
func run() error {
	cfg := platform.LoadConfig()

	// --- Adaptateurs outbound (pilotés) ---
	db, err := sqlite.Open(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("ouverture de la base SQLite (%s) : %w", cfg.DBPath, err)
	}
	defer func() { _ = db.Close() }()

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
		DevMode:       cfg.DevMode,
		DevPassword:   cfg.DevPassword,
	})
	taskService := service.NewTaskService(taskRepo, time.Now, uuid.NewString)

	// --- Bootstrap : admin initial (mot de passe affiché une fois dans les logs) ---
	res, err := authService.EnsureAdmin(context.Background())
	if err != nil {
		return fmt.Errorf("initialisation de l'admin : %w", err)
	}
	// GeneratedPassword est renseigné à la création (prod : une fois) et à chaque
	// démarrage en dev (mot de passe réinitialisé/réaffiché).
	if res.GeneratedPassword != "" {
		printAdminBanner(res.Username, res.GeneratedPassword, res.Created)
	}

	// --- Adaptateur inbound (pilotant) ---
	router := adapterhttp.NewRouter(
		cfg.CORSOrigins,
		adapterhttp.CookieConfig{Name: cfg.CookieName, Secure: cfg.CookieSecure, TTL: cfg.SessionTTL},
		authService,
		taskService,
	)

	slog.Info("démarrage de l'API Graefik", "port", cfg.Port, "env", cfg.Env)
	return router.Start(":" + cfg.Port)
}

// printAdminBanner affiche le compte admin dans les logs.
// created distingue la création initiale (prod, une fois) de la réinitialisation
// dev (réaffichée à chaque démarrage).
func printAdminBanner(username, password string, created bool) {
	const line = "============================================================"
	fmt.Println("\n" + line)
	if created {
		fmt.Println("  Graefik — compte administrateur initial")
	} else {
		fmt.Println("  Graefik — compte administrateur (mode dev)")
	}
	fmt.Println()
	fmt.Printf("  Identifiant  : %s\n", username)
	fmt.Printf("  Mot de passe : %s\n", password)
	fmt.Println()
	if created {
		fmt.Println("  Ce mot de passe n'est affiché qu'UNE seule fois.")
		fmt.Println("  Notez-le et changez-le dès que possible.")
	} else {
		fmt.Println("  Mode dev : mot de passe réinitialisé et réaffiché à chaque démarrage.")
	}
	fmt.Printf("%s\n\n", line)
}

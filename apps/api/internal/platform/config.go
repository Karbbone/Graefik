// Package platform regroupe les préoccupations techniques transverses :
// configuration, génération de secrets, (futur) logger, bootstrap serveur.
package platform

import (
	"os"
	"strings"
	"time"
)

// Config regroupe la configuration de l'application, chargée depuis l'environnement.
type Config struct {
	Port         string
	Env          string
	CORSOrigins  []string
	DBPath       string
	CookieName   string
	CookieSecure string // "auto" | "true" | "false"
	SessionTTL   time.Duration
}

// LoadConfig lit la configuration depuis les variables d'environnement.
func LoadConfig() Config {
	return Config{
		Port:         getenv("PORT", "8080"),
		Env:          os.Getenv("APP_ENV"),
		CORSOrigins:  parseOrigins(os.Getenv("CORS_ORIGINS")),
		DBPath:       getenv("DB_PATH", "/data/graefik.db"),
		CookieName:   getenv("COOKIE_NAME", "graefik_session"),
		CookieSecure: getenv("COOKIE_SECURE", "auto"),
		SessionTTL:   getduration("SESSION_TTL", 7*24*time.Hour),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getduration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

// parseOrigins découpe une liste d'origines séparées par des virgules.
// Retourne ["*"] si rien n'est défini.
func parseOrigins(raw string) []string {
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

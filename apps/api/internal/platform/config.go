// Package platform regroupe les préoccupations techniques transverses :
// configuration, logger, bootstrap serveur, (futur) pool de base de données.
package platform

import (
	"os"
	"strings"
)

// Config regroupe la configuration de l'application, chargée depuis l'environnement.
type Config struct {
	Port        string
	Env         string
	CORSOrigins []string
}

// LoadConfig lit la configuration depuis les variables d'environnement.
func LoadConfig() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		Env:         os.Getenv("APP_ENV"),
		CORSOrigins: parseOrigins(os.Getenv("CORS_ORIGINS")),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
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

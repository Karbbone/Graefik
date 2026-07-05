package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

// CookieConfig configure le cookie de session.
type CookieConfig struct {
	Name   string
	Secure string // "auto" | "true" | "false"
	TTL    time.Duration
}

// determineSecure décide de la valeur du flag Secure selon la config et la requête.
func determineSecure(c *echo.Context, mode string) bool {
	switch strings.ToLower(mode) {
	case "true":
		return true
	case "false":
		return false
	default: // "auto"
		if c.Request().TLS != nil {
			return true
		}
		return strings.EqualFold(c.Request().Header.Get("X-Forwarded-Proto"), "https")
	}
}

func buildSessionCookie(name, value string, expires time.Time, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

func clearSessionCookie(name string, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
}

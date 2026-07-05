package http

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// userContextKey est la clé sous laquelle l'utilisateur authentifié est stocké dans le contexte Echo.
const userContextKey = "user"

// RequireAuth est un middleware qui exige une session valide (cookie).
// Il place l'utilisateur authentifié dans le contexte sous la clé "user".
func RequireAuth(authenticate port.AuthenticateUseCase, cookieName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			cookie, err := c.Request().Cookie(cookieName)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "non authentifié")
			}
			user, err := authenticate.Execute(c.Request().Context(), cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "session invalide")
			}
			c.Set(userContextKey, user)
			return next(c)
		}
	}
}

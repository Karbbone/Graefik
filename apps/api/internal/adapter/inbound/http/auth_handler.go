package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// AuthUseCases regroupe les cas d'usage d'authentification consommés par l'adaptateur HTTP.
type AuthUseCases struct {
	Login        port.LoginUseCase
	Logout       port.LogoutUseCase
	Authenticate port.AuthenticateUseCase
}

// AuthHandler expose les endpoints d'authentification.
type AuthHandler struct {
	uc      AuthUseCases
	cookie  CookieConfig
	limiter *loginLimiter
}

// NewAuthHandler construit le handler (verrou anti-brute-force : 5 échecs -> 1 min).
func NewAuthHandler(uc AuthUseCases, cookie CookieConfig) *AuthHandler {
	return &AuthHandler{
		uc:      uc,
		cookie:  cookie,
		limiter: newLoginLimiter(5, time.Minute, time.Now),
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type userResponse struct {
	Username string `json:"username"`
}

// Login — POST /api/auth/login
func (h *AuthHandler) Login(c *echo.Context) error {
	ip := c.RealIP()
	if !h.limiter.allowed(ip) {
		return echo.NewHTTPError(http.StatusTooManyRequests, "trop de tentatives, réessayez plus tard")
	}

	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "corps de requête invalide")
	}

	session, err := h.uc.Login.Execute(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			h.limiter.recordFailure(ip)
			return echo.NewHTTPError(http.StatusUnauthorized, "identifiants invalides")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "erreur interne")
	}
	h.limiter.reset(ip)

	secure := determineSecure(c, h.cookie.Secure)
	c.SetCookie(buildSessionCookie(h.cookie.Name, session.Token, session.ExpiresAt, secure))

	return c.JSON(http.StatusOK, userResponse{Username: req.Username})
}

// Logout — POST /api/auth/logout
func (h *AuthHandler) Logout(c *echo.Context) error {
	if cookie, err := c.Request().Cookie(h.cookie.Name); err == nil {
		_ = h.uc.Logout.Execute(c.Request().Context(), cookie.Value)
	}
	secure := determineSecure(c, h.cookie.Secure)
	c.SetCookie(clearSessionCookie(h.cookie.Name, secure))
	return c.NoContent(http.StatusNoContent)
}

// Me — GET /api/auth/me
func (h *AuthHandler) Me(c *echo.Context) error {
	cookie, err := c.Request().Cookie(h.cookie.Name)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "non authentifié")
	}
	user, err := h.uc.Authenticate.Execute(c.Request().Context(), cookie.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "session invalide")
	}
	return c.JSON(http.StatusOK, userResponse{Username: user.Username})
}

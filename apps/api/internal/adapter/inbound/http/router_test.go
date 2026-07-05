package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	adapterhttp "github.com/Karbbone/Graefik/apps/api/internal/adapter/inbound/http"
	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
)

const testCookieName = "graefik_session"

// testUseCases fournit un jeu de mocks pour les trois use cases d'auth.
type testUseCases struct {
	login        *portmocks.MockLoginUseCase
	logout       *portmocks.MockLogoutUseCase
	authenticate *portmocks.MockAuthenticateUseCase
}

func newTestUseCases(t *testing.T) testUseCases {
	return testUseCases{
		login:        portmocks.NewMockLoginUseCase(t),
		logout:       portmocks.NewMockLogoutUseCase(t),
		authenticate: portmocks.NewMockAuthenticateUseCase(t),
	}
}

func (uc testUseCases) router() *echo.Echo {
	return adapterhttp.NewRouter(
		[]string{"*"},
		adapterhttp.CookieConfig{Name: testCookieName, Secure: "false", TTL: time.Hour},
		adapterhttp.AuthUseCases{Login: uc.login, Logout: uc.logout, Authenticate: uc.authenticate},
	)
}

func TestRouter_Health_Public(t *testing.T) {
	e := newTestUseCases(t).router()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestRouter_Login_OK(t *testing.T) {
	uc := newTestUseCases(t)
	session := &domain.Session{Token: "tok-123", UserID: "u1", ExpiresAt: time.Now().Add(time.Hour)}
	uc.login.EXPECT().Execute(mock.Anything, "graefik", "secret").Return(session, nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"graefik","password":"secret"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	uc.router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"username":"graefik"`)
	assert.Contains(t, rec.Header().Get("Set-Cookie"), testCookieName+"=tok-123")
	assert.Contains(t, rec.Header().Get("Set-Cookie"), "HttpOnly")
}

func TestRouter_Login_BadCredentials(t *testing.T) {
	uc := newTestUseCases(t)
	uc.login.EXPECT().Execute(mock.Anything, "graefik", "wrong").
		Return(nil, domain.ErrInvalidCredentials).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"graefik","password":"wrong"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	uc.router().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Empty(t, rec.Header().Get("Set-Cookie"))
}

func TestRouter_Me_NoCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	newTestUseCases(t).router().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRouter_Me_ValidSession(t *testing.T) {
	uc := newTestUseCases(t)
	uc.authenticate.EXPECT().Execute(mock.Anything, "tok-123").
		Return(&domain.User{ID: "u1", Username: "graefik"}, nil).Once()

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: "tok-123"})
	rec := httptest.NewRecorder()
	uc.router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"username":"graefik"`)
}

func TestRouter_Logout_ClearsCookie(t *testing.T) {
	uc := newTestUseCases(t)
	uc.logout.EXPECT().Execute(mock.Anything, "tok-123").Return(nil).Once()

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: "tok-123"})
	rec := httptest.NewRecorder()
	uc.router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	assert.Contains(t, rec.Header().Get("Set-Cookie"), testCookieName+"=")
}

// TestRequireAuth couvre le middleware de protection sur une route factice.
func TestRequireAuth(t *testing.T) {
	authenticate := portmocks.NewMockAuthenticateUseCase(t)
	authenticate.EXPECT().Execute(mock.Anything, "tok-123").
		Return(&domain.User{ID: "u1", Username: "graefik"}, nil).Once()

	e := echo.New()
	e.GET("/protected",
		func(c *echo.Context) error { return c.NoContent(http.StatusOK) },
		adapterhttp.RequireAuth(authenticate, testCookieName),
	)

	// Sans cookie → 401.
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Avec cookie valide → 200.
	req = httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: "tok-123"})
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

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

func newTestRouter(auth *portmocks.MockAuthService, tasks *portmocks.MockTaskService) *echo.Echo {
	return adapterhttp.NewRouter(
		[]string{"*"},
		adapterhttp.CookieConfig{Name: testCookieName, Secure: "false", TTL: time.Hour},
		auth,
		tasks,
	)
}

func TestRouter_Health_Public(t *testing.T) {
	e := newTestRouter(portmocks.NewMockAuthService(t), portmocks.NewMockTaskService(t))

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

func TestRouter_Login_OK(t *testing.T) {
	auth := portmocks.NewMockAuthService(t)
	session := &domain.Session{Token: "tok-123", UserID: "u1", ExpiresAt: time.Now().Add(time.Hour)}
	auth.EXPECT().Login(mock.Anything, "graefik", "secret").Return(session, nil).Once()

	e := newTestRouter(auth, portmocks.NewMockTaskService(t))

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"graefik","password":"secret"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"username":"graefik"`)
	// Le cookie de session est posé.
	assert.Contains(t, rec.Header().Get("Set-Cookie"), testCookieName+"=tok-123")
	assert.Contains(t, rec.Header().Get("Set-Cookie"), "HttpOnly")
}

func TestRouter_Login_BadCredentials(t *testing.T) {
	auth := portmocks.NewMockAuthService(t)
	auth.EXPECT().Login(mock.Anything, "graefik", "wrong").
		Return(nil, domain.ErrInvalidCredentials).Once()

	e := newTestRouter(auth, portmocks.NewMockTaskService(t))

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login",
		strings.NewReader(`{"username":"graefik","password":"wrong"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Empty(t, rec.Header().Get("Set-Cookie"))
}

func TestRouter_Me_NoCookie(t *testing.T) {
	e := newTestRouter(portmocks.NewMockAuthService(t), portmocks.NewMockTaskService(t))

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRouter_Me_ValidSession(t *testing.T) {
	auth := portmocks.NewMockAuthService(t)
	auth.EXPECT().Authenticate(mock.Anything, "tok-123").
		Return(&domain.User{ID: "u1", Username: "graefik"}, nil).Once()

	e := newTestRouter(auth, portmocks.NewMockTaskService(t))

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: "tok-123"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"username":"graefik"`)
}

func TestRouter_Tasks_RequiresAuth(t *testing.T) {
	e := newTestRouter(portmocks.NewMockAuthService(t), portmocks.NewMockTaskService(t))

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRouter_Tasks_List_Authenticated(t *testing.T) {
	auth := portmocks.NewMockAuthService(t)
	auth.EXPECT().Authenticate(mock.Anything, "tok-123").
		Return(&domain.User{ID: "u1", Username: "graefik"}, nil).Once()

	tasks := portmocks.NewMockTaskService(t)
	tasks.EXPECT().List(mock.Anything).Return([]*domain.Task{}, nil).Once()

	e := newTestRouter(auth, tasks)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: "tok-123"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

func TestRouter_Logout_ClearsCookie(t *testing.T) {
	auth := portmocks.NewMockAuthService(t)
	auth.EXPECT().Logout(mock.Anything, "tok-123").Return(nil).Once()

	e := newTestRouter(auth, portmocks.NewMockTaskService(t))

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	req.AddCookie(&http.Cookie{Name: testCookieName, Value: "tok-123"})
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	// Cookie effacé (Max-Age=0 ou expiration passée).
	assert.Contains(t, rec.Header().Get("Set-Cookie"), testCookieName+"=")
}

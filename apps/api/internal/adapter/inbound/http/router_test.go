package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	adapterhttp "github.com/Karbbone/Graefik/apps/api/internal/adapter/inbound/http"
	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRouter_CreateTask(t *testing.T) {
	created, _ := domain.NewTask("1", "Nouvelle tâche", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	svc := portmocks.NewMockTaskService(t)
	svc.EXPECT().Create(mock.Anything, "Nouvelle tâche").Return(created, nil).Once()

	e := adapterhttp.NewRouter([]string{"*"}, svc)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"title":"Nouvelle tâche"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "Nouvelle tâche")
}

func TestRouter_CreateTask_TitreVide(t *testing.T) {
	svc := portmocks.NewMockTaskService(t)
	svc.EXPECT().
		Create(mock.Anything, "").
		Return(nil, domain.ErrTaskTitleRequired).
		Once()

	e := adapterhttp.NewRouter([]string{"*"}, svc)

	req := httptest.NewRequest(http.MethodPost, "/api/tasks", strings.NewReader(`{"title":""}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRouter_ListTasks(t *testing.T) {
	svc := portmocks.NewMockTaskService(t)
	svc.EXPECT().List(mock.Anything).Return([]*domain.Task{}, nil).Once()

	e := adapterhttp.NewRouter([]string{"*"}, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

func TestRouter_Health(t *testing.T) {
	svc := portmocks.NewMockTaskService(t)
	e := adapterhttp.NewRouter([]string{"*"}, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"status":"ok"`)
}

package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/service"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// deps fournit une horloge et un générateur d'ID déterministes pour les tests.
func deps() (func() time.Time, func() string) {
	now := func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }
	newID := func() string { return "task-1" }
	return now, newID
}

func TestTaskService_Create_OK(t *testing.T) {
	repo := portmocks.NewMockTaskRepository(t)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil).Once()

	now, newID := deps()
	svc := service.NewTaskService(repo, now, newID)

	task, err := svc.Create(context.Background(), "Ma tâche")

	require.NoError(t, err)
	assert.Equal(t, "task-1", task.ID)
	assert.Equal(t, "Ma tâche", task.Title)
	assert.False(t, task.Done)
	// Les attentes du mock sont vérifiées automatiquement au cleanup (via t).
}

func TestTaskService_Create_TitreVide(t *testing.T) {
	repo := portmocks.NewMockTaskRepository(t)
	// Aucune attente Save : le repository ne doit PAS être appelé si le titre est invalide.

	now, newID := deps()
	svc := service.NewTaskService(repo, now, newID)

	_, err := svc.Create(context.Background(), "   ")

	assert.ErrorIs(t, err, domain.ErrTaskTitleRequired)
}

func TestTaskService_Create_ErreurRepo(t *testing.T) {
	repoErr := errors.New("base de données indisponible")
	repo := portmocks.NewMockTaskRepository(t)
	repo.EXPECT().Save(mock.Anything, mock.Anything).Return(repoErr).Once()

	now, newID := deps()
	svc := service.NewTaskService(repo, now, newID)

	_, err := svc.Create(context.Background(), "Tâche valide")

	assert.ErrorIs(t, err, repoErr)
}

func TestTaskService_List(t *testing.T) {
	existing := []*domain.Task{{ID: "1", Title: "A"}, {ID: "2", Title: "B"}}
	repo := portmocks.NewMockTaskRepository(t)
	repo.EXPECT().FindAll(mock.Anything).Return(existing, nil).Once()

	now, newID := deps()
	svc := service.NewTaskService(repo, now, newID)

	tasks, err := svc.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, tasks, 2)
}

package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/adapter/outbound/memory"
	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_SaveAndFindAll(t *testing.T) {
	repo := memory.NewTaskRepository()
	ctx := context.Background()

	// Vide au départ.
	all, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, all)

	task, err := domain.NewTask("1", "Première tâche", time.Now())
	require.NoError(t, err)
	require.NoError(t, repo.Save(ctx, task))

	all, err = repo.FindAll(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)
	assert.Equal(t, "1", all[0].ID)
}

func TestTaskRepository_FindAllRetourneUneCopie(t *testing.T) {
	repo := memory.NewTaskRepository()
	ctx := context.Background()

	task, _ := domain.NewTask("1", "Tâche", time.Now())
	require.NoError(t, repo.Save(ctx, task))

	first, _ := repo.FindAll(ctx)
	first[0] = nil // muter la tranche retournée ne doit pas affecter le repository

	second, _ := repo.FindAll(ctx)
	require.Len(t, second, 1)
	assert.NotNil(t, second[0])
}

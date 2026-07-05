package domain_test

import (
	"testing"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	createdAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		title     string
		wantTitle string
		wantErr   error
	}{
		{name: "titre valide", title: "Acheter du pain", wantTitle: "Acheter du pain"},
		{name: "titre entouré d'espaces", title: "  Ranger le bureau  ", wantTitle: "Ranger le bureau"},
		{name: "titre vide", title: "   ", wantErr: domain.ErrTaskTitleRequired},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task, err := domain.NewTask("id-1", tc.title, createdAt)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantTitle, task.Title)
			assert.False(t, task.Done)
			assert.Equal(t, createdAt, task.CreatedAt)
		})
	}
}

func TestTask_MarkDone(t *testing.T) {
	task, err := domain.NewTask("id-1", "Tâche", time.Now())
	require.NoError(t, err)

	task.MarkDone()

	assert.True(t, task.Done)
}

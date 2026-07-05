package port

import (
	"context"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
)

// TaskRepository est un port outbound : la persistance des tâches.
// Implémenté par un adaptateur piloté (en mémoire aujourd'hui, Postgres demain).
type TaskRepository interface {
	Save(ctx context.Context, task *domain.Task) error
	FindAll(ctx context.Context) ([]*domain.Task, error)
}

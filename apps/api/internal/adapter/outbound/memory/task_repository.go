// Package memory est un adaptateur OUTBOUND (piloté) : une implémentation
// en mémoire, thread-safe, du port port.TaskRepository. Utile en dev et en test ;
// un adaptateur Postgres viendra le remplacer sans toucher au cœur.
package memory

import (
	"context"
	"sync"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// TaskRepository stocke les tâches en mémoire.
type TaskRepository struct {
	mu    sync.RWMutex
	tasks []*domain.Task
}

// Vérification à la compilation que l'adaptateur satisfait le port outbound.
var _ port.TaskRepository = (*TaskRepository)(nil)

// NewTaskRepository crée un repository en mémoire vide.
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// Save ajoute une tâche.
func (r *TaskRepository) Save(_ context.Context, task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = append(r.tasks, task)
	return nil
}

// FindAll retourne une copie de toutes les tâches.
func (r *TaskRepository) FindAll(_ context.Context) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*domain.Task, len(r.tasks))
	copy(out, r.tasks)
	return out, nil
}

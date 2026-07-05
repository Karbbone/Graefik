// Package service contient les use cases : ils implémentent les ports inbound
// et orchestrent le domaine + les ports outbound. Aucune dépendance framework.
package service

import (
	"context"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// TaskService implémente le port inbound port.TaskService.
type TaskService struct {
	repo  port.TaskRepository
	now   func() time.Time
	newID func() string
}

// Vérification à la compilation que TaskService satisfait le port inbound.
var _ port.TaskService = (*TaskService)(nil)

// NewTaskService injecte les dépendances : le repository (port outbound),
// une horloge et un générateur d'identifiants (injectés pour la testabilité).
func NewTaskService(repo port.TaskRepository, now func() time.Time, newID func() string) *TaskService {
	return &TaskService{repo: repo, now: now, newID: newID}
}

// Create valide puis persiste une nouvelle tâche.
func (s *TaskService) Create(ctx context.Context, title string) (*domain.Task, error) {
	task, err := domain.NewTask(s.newID(), title, s.now())
	if err != nil {
		return nil, err
	}
	if err := s.repo.Save(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

// List retourne toutes les tâches.
func (s *TaskService) List(ctx context.Context) ([]*domain.Task, error) {
	return s.repo.FindAll(ctx)
}

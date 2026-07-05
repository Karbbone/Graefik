// Package domain contient le cœur métier : entités, value objects, règles et
// erreurs. Il ne dépend d'AUCUN framework ni détail technique (ni Echo, ni BDD).
package domain

import (
	"errors"
	"strings"
	"time"
)

// Erreurs métier du domaine Task.
var (
	ErrTaskTitleRequired = errors.New("le titre de la tâche est requis")
	ErrTaskNotFound      = errors.New("tâche introuvable")
)

// Task est l'entité racine (aggregate root) du domaine des tâches.
type Task struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt time.Time
}

// NewTask crée une tâche valide en appliquant les invariants métier.
// La création échoue si le titre est vide.
func NewTask(id, title string, createdAt time.Time) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrTaskTitleRequired
	}
	return &Task{
		ID:        id,
		Title:     title,
		Done:      false,
		CreatedAt: createdAt,
	}, nil
}

// MarkDone marque la tâche comme terminée.
func (t *Task) MarkDone() {
	t.Done = true
}

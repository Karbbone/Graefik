// Package port déclare les interfaces (ports) de l'hexagone.
//
//   - Ports INBOUND (driving) : cas d'usage exposés par l'application, appelés
//     par les adaptateurs pilotants (HTTP, CLI…).
//   - Ports OUTBOUND (driven) : dépendances dont l'application a besoin,
//     implémentées par les adaptateurs pilotés (persistance, clients externes…).
package port

import (
	"context"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
)

// TaskService est un port inbound : les cas d'usage du domaine des tâches.
type TaskService interface {
	Create(ctx context.Context, title string) (*domain.Task, error)
	List(ctx context.Context) ([]*domain.Task, error)
}

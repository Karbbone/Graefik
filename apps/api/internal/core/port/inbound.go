// Package port déclare les interfaces (ports) de l'hexagone.
//
//   - Ports INBOUND (driving) : les cas d'usage exposés par l'application, un par
//     opération, appelés par les adaptateurs pilotants (HTTP, CLI…).
//   - Ports OUTBOUND (driven) : les dépendances dont les use cases ont besoin,
//     implémentées par les adaptateurs pilotés (persistance, hachage…).
package port

import (
	"context"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
)

// LoginUseCase : authentifie et ouvre une session.
type LoginUseCase interface {
	Execute(ctx context.Context, username, password string) (*domain.Session, error)
}

// LogoutUseCase : invalide une session.
type LogoutUseCase interface {
	Execute(ctx context.Context, token string) error
}

// AuthenticateUseCase : valide un token de session et renvoie l'utilisateur.
type AuthenticateUseCase interface {
	Execute(ctx context.Context, token string) (*domain.User, error)
}

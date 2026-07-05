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

// BootstrapResult décrit le résultat de la création éventuelle de l'admin initial.
type BootstrapResult struct {
	Created           bool
	Username          string
	GeneratedPassword string // renseigné uniquement si Created == true
}

// AuthService est un port inbound : les cas d'usage d'authentification.
type AuthService interface {
	// EnsureAdmin crée l'utilisateur admin initial s'il n'existe aucun compte.
	EnsureAdmin(ctx context.Context) (BootstrapResult, error)
	// Login vérifie les identifiants et ouvre une session.
	Login(ctx context.Context, username, password string) (*domain.Session, error)
	// Logout invalide la session correspondant au token.
	Logout(ctx context.Context, token string) error
	// Authenticate valide un token de session et renvoie l'utilisateur (prolonge la session).
	Authenticate(ctx context.Context, token string) (*domain.User, error)
}

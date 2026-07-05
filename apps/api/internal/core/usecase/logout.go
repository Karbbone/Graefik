package usecase

import (
	"context"

	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// Logout invalide une session.
type Logout struct {
	sessions port.SessionRepository
}

var _ port.LogoutUseCase = (*Logout)(nil)

// NewLogout construit le cas d'usage Logout.
func NewLogout(sessions port.SessionRepository) *Logout {
	return &Logout{sessions: sessions}
}

// Execute supprime la session (idempotent).
func (uc *Logout) Execute(ctx context.Context, token string) error {
	return uc.sessions.Delete(ctx, token)
}

package usecase

import (
	"context"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// Authenticate valide un token de session et renvoie l'utilisateur associé.
type Authenticate struct {
	users    port.UserRepository
	sessions port.SessionRepository
	now      func() time.Time
	ttl      time.Duration
}

var _ port.AuthenticateUseCase = (*Authenticate)(nil)

// NewAuthenticate construit le cas d'usage Authenticate.
func NewAuthenticate(
	users port.UserRepository,
	sessions port.SessionRepository,
	now func() time.Time,
	ttl time.Duration,
) *Authenticate {
	return &Authenticate{users: users, sessions: sessions, now: now, ttl: ttl}
}

// Execute valide le token, prolonge la session (glissante) et renvoie l'utilisateur.
func (uc *Authenticate) Execute(ctx context.Context, token string) (*domain.User, error) {
	session, err := uc.sessions.FindByToken(ctx, token)
	if err != nil {
		return nil, domain.ErrSessionInvalid
	}

	now := uc.now()
	if session.IsExpired(now) {
		_ = uc.sessions.Delete(ctx, token)
		return nil, domain.ErrSessionInvalid
	}

	// Session glissante : on repousse l'expiration (best-effort).
	_ = uc.sessions.UpdateExpiry(ctx, token, now.Add(uc.ttl))

	user, err := uc.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, domain.ErrSessionInvalid
	}
	return user, nil
}

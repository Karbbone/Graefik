// Package usecase contient les cas d'usage de l'application : un cas d'usage =
// une opération métier, dans son propre fichier. Chaque use case implémente un
// port inbound et ne dépend que des ports outbound dont il a réellement besoin.
package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// Login authentifie un utilisateur et ouvre une session.
type Login struct {
	users    port.UserRepository
	sessions port.SessionRepository
	hasher   port.PasswordHasher
	now      func() time.Time
	newToken func() (string, error)
	ttl      time.Duration
}

var _ port.LoginUseCase = (*Login)(nil)

// NewLogin construit le cas d'usage Login.
func NewLogin(
	users port.UserRepository,
	sessions port.SessionRepository,
	hasher port.PasswordHasher,
	now func() time.Time,
	newToken func() (string, error),
	ttl time.Duration,
) *Login {
	return &Login{users: users, sessions: sessions, hasher: hasher, now: now, newToken: newToken, ttl: ttl}
}

// Execute vérifie les identifiants puis crée une session.
func (uc *Login) Execute(ctx context.Context, username, password string) (*domain.Session, error) {
	user, err := uc.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// Erreur générique : on ne divulgue pas quel champ est faux.
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := uc.hasher.Compare(user.PasswordHash, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := uc.newToken()
	if err != nil {
		return nil, err
	}
	now := uc.now()
	session := &domain.Session{
		Token:     token,
		UserID:    user.ID,
		CreatedAt: now,
		ExpiresAt: now.Add(uc.ttl),
	}
	if err := uc.sessions.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

package port

import (
	"context"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
)

// TaskRepository est un port outbound : la persistance des tâches.
// Implémenté par un adaptateur piloté (en mémoire aujourd'hui, Postgres demain).
type TaskRepository interface {
	Save(ctx context.Context, task *domain.Task) error
	FindAll(ctx context.Context) ([]*domain.Task, error)
}

// UserRepository est un port outbound : la persistance des utilisateurs.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdatePassword(ctx context.Context, id, passwordHash string) error
	Count(ctx context.Context) (int, error)
}

// SessionRepository est un port outbound : la persistance des sessions.
type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindByToken(ctx context.Context, token string) (*domain.Session, error)
	UpdateExpiry(ctx context.Context, token string, expiresAt time.Time) error
	Delete(ctx context.Context, token string) error
	DeleteExpired(ctx context.Context, now time.Time) error
}

// PasswordHasher est un port outbound : le hachage/vérification des mots de passe.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

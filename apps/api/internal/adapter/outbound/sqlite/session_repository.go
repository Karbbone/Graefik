package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// SessionRepository persiste les sessions dans SQLite.
type SessionRepository struct {
	db *sql.DB
}

var _ port.SessionRepository = (*SessionRepository)(nil)

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, s *domain.Session) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO sessions (token, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)`,
		s.Token, s.UserID, s.CreatedAt.Unix(), s.ExpiresAt.Unix(),
	)
	return err
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*domain.Session, error) {
	var (
		s         domain.Session
		createdAt int64
		expiresAt int64
	)
	err := r.db.QueryRowContext(ctx,
		`SELECT token, user_id, created_at, expires_at FROM sessions WHERE token = ?`, token).
		Scan(&s.Token, &s.UserID, &createdAt, &expiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrSessionInvalid
	}
	if err != nil {
		return nil, err
	}
	s.CreatedAt = time.Unix(createdAt, 0).UTC()
	s.ExpiresAt = time.Unix(expiresAt, 0).UTC()
	return &s, nil
}

func (r *SessionRepository) UpdateExpiry(ctx context.Context, token string, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET expires_at = ? WHERE token = ?`, expiresAt.Unix(), token)
	return err
}

func (r *SessionRepository) Delete(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE token = ?`, token)
	return err
}

func (r *SessionRepository) DeleteExpired(ctx context.Context, now time.Time) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at < ?`, now.Unix())
	return err
}

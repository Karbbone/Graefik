package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// UserRepository persiste les utilisateurs dans SQLite.
type UserRepository struct {
	db *sql.DB
}

var _ port.UserRepository = (*UserRepository)(nil)

// NewUserRepository crée un repository d'utilisateurs sur la base fournie.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create insère un nouvel utilisateur.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, username, password_hash, created_at) VALUES (?, ?, ?, ?)`,
		user.ID, user.Username, user.PasswordHash, user.CreatedAt.Unix(),
	)
	return err
}

// FindByID retourne l'utilisateur d'identifiant donné (ErrUserNotFound sinon).
func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	return r.queryUser(ctx, `SELECT id, username, password_hash, created_at FROM users WHERE id = ?`, id)
}

// FindByUsername retourne l'utilisateur au nom donné (ErrUserNotFound sinon).
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	return r.queryUser(ctx, `SELECT id, username, password_hash, created_at FROM users WHERE username = ?`, username)
}

// UpdatePassword remplace le hash du mot de passe d'un utilisateur.
func (r *UserRepository) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, id)
	return err
}

// Count retourne le nombre d'utilisateurs.
func (r *UserRepository) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&n)
	return n, err
}

func (r *UserRepository) queryUser(ctx context.Context, query string, arg any) (*domain.User, error) {
	var (
		u         domain.User
		createdAt int64
	)
	err := r.db.QueryRowContext(ctx, query, arg).
		Scan(&u.ID, &u.Username, &u.PasswordHash, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	u.CreatedAt = time.Unix(createdAt, 0).UTC()
	return &u, nil
}

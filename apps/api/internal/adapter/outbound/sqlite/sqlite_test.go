package sqlite_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Karbbone/Graefik/apps/api/internal/adapter/outbound/sqlite"
	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
)

func newDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sqlite.Open(":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestUserRepository(t *testing.T) {
	repo := sqlite.NewUserRepository(newDB(t))
	ctx := context.Background()

	n, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, n)

	_, err = repo.FindByUsername(ctx, "graefik")
	assert.ErrorIs(t, err, domain.ErrUserNotFound)

	user := &domain.User{ID: "1", Username: "graefik", PasswordHash: "hash", CreatedAt: time.Now()}
	require.NoError(t, repo.Create(ctx, user))

	byName, err := repo.FindByUsername(ctx, "graefik")
	require.NoError(t, err)
	assert.Equal(t, "1", byName.ID)
	assert.Equal(t, "hash", byName.PasswordHash)

	byID, err := repo.FindByID(ctx, "1")
	require.NoError(t, err)
	assert.Equal(t, "graefik", byID.Username)

	n, err = repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, n)
}

func TestSessionRepository(t *testing.T) {
	repo := sqlite.NewSessionRepository(newDB(t))
	ctx := context.Background()

	_, err := repo.FindByToken(ctx, "nope")
	assert.ErrorIs(t, err, domain.ErrSessionInvalid)

	now := time.Now().Truncate(time.Second)
	s := &domain.Session{Token: "t1", UserID: "u1", CreatedAt: now, ExpiresAt: now.Add(time.Hour)}
	require.NoError(t, repo.Create(ctx, s))

	got, err := repo.FindByToken(ctx, "t1")
	require.NoError(t, err)
	assert.Equal(t, "u1", got.UserID)

	require.NoError(t, repo.UpdateExpiry(ctx, "t1", now.Add(2*time.Hour)))
	got, err = repo.FindByToken(ctx, "t1")
	require.NoError(t, err)
	assert.Equal(t, now.Add(2*time.Hour).Unix(), got.ExpiresAt.Unix())

	// Nettoyage des sessions expirées.
	old := &domain.Session{Token: "t2", UserID: "u1", CreatedAt: now.Add(-2 * time.Hour), ExpiresAt: now.Add(-time.Hour)}
	require.NoError(t, repo.Create(ctx, old))
	require.NoError(t, repo.DeleteExpired(ctx, now))
	_, err = repo.FindByToken(ctx, "t2")
	assert.ErrorIs(t, err, domain.ErrSessionInvalid)

	require.NoError(t, repo.Delete(ctx, "t1"))
	_, err = repo.FindByToken(ctx, "t1")
	assert.ErrorIs(t, err, domain.ErrSessionInvalid)
}

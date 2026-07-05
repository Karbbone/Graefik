package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/usecase"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
)

func fixedNow() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }

func newToken() (string, error) { return "tok-1", nil }

func TestLogin_OK(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	sessions := portmocks.NewMockSessionRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	user := &domain.User{ID: "u1", Username: "graefik", PasswordHash: "hashed"}
	users.EXPECT().FindByUsername(mock.Anything, "graefik").Return(user, nil).Once()
	hasher.EXPECT().Compare("hashed", "secret").Return(nil).Once()
	sessions.EXPECT().Create(mock.Anything, mock.MatchedBy(func(s *domain.Session) bool {
		return s.Token == "tok-1" && s.UserID == "u1"
	})).Return(nil).Once()

	uc := usecase.NewLogin(users, sessions, hasher, fixedNow, newToken, time.Hour)

	session, err := uc.Execute(context.Background(), "graefik", "secret")

	require.NoError(t, err)
	assert.Equal(t, "tok-1", session.Token)
	assert.Equal(t, "u1", session.UserID)
}

func TestLogin_UtilisateurInconnu(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	users.EXPECT().FindByUsername(mock.Anything, "ghost").Return(nil, domain.ErrUserNotFound).Once()

	uc := usecase.NewLogin(users, portmocks.NewMockSessionRepository(t), portmocks.NewMockPasswordHasher(t), fixedNow, newToken, time.Hour)

	_, err := uc.Execute(context.Background(), "ghost", "x")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestLogin_MauvaisMotDePasse(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	user := &domain.User{ID: "u1", Username: "graefik", PasswordHash: "hashed"}
	users.EXPECT().FindByUsername(mock.Anything, "graefik").Return(user, nil).Once()
	hasher.EXPECT().Compare("hashed", "bad").Return(errors.New("mismatch")).Once()

	uc := usecase.NewLogin(users, portmocks.NewMockSessionRepository(t), hasher, fixedNow, newToken, time.Hour)

	_, err := uc.Execute(context.Background(), "graefik", "bad")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

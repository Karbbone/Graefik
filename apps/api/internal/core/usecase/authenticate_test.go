package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/usecase"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
)

func TestAuthenticate_Valide(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	sessions := portmocks.NewMockSessionRepository(t)
	now := fixedNow()

	session := &domain.Session{Token: "tok-1", UserID: "u1", ExpiresAt: now.Add(time.Hour)}
	sessions.EXPECT().FindByToken(mock.Anything, "tok-1").Return(session, nil).Once()
	sessions.EXPECT().UpdateExpiry(mock.Anything, "tok-1", mock.Anything).Return(nil).Once()
	users.EXPECT().FindByID(mock.Anything, "u1").Return(&domain.User{ID: "u1", Username: "graefik"}, nil).Once()

	uc := usecase.NewAuthenticate(users, sessions, fixedNow, time.Hour)

	user, err := uc.Execute(context.Background(), "tok-1")

	require.NoError(t, err)
	assert.Equal(t, "graefik", user.Username)
}

func TestAuthenticate_Expiree(t *testing.T) {
	sessions := portmocks.NewMockSessionRepository(t)
	now := fixedNow()

	expired := &domain.Session{Token: "tok-1", UserID: "u1", ExpiresAt: now.Add(-time.Minute)}
	sessions.EXPECT().FindByToken(mock.Anything, "tok-1").Return(expired, nil).Once()
	sessions.EXPECT().Delete(mock.Anything, "tok-1").Return(nil).Once()

	uc := usecase.NewAuthenticate(portmocks.NewMockUserRepository(t), sessions, fixedNow, time.Hour)

	_, err := uc.Execute(context.Background(), "tok-1")

	assert.ErrorIs(t, err, domain.ErrSessionInvalid)
}

func TestAuthenticate_Inexistante(t *testing.T) {
	sessions := portmocks.NewMockSessionRepository(t)
	sessions.EXPECT().FindByToken(mock.Anything, "nope").Return(nil, domain.ErrSessionInvalid).Once()

	uc := usecase.NewAuthenticate(portmocks.NewMockUserRepository(t), sessions, fixedNow, time.Hour)

	_, err := uc.Execute(context.Background(), "nope")

	assert.ErrorIs(t, err, domain.ErrSessionInvalid)
}

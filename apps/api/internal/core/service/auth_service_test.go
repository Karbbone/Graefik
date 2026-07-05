package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/service"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
)

func authCfg() service.AuthConfig {
	return service.AuthConfig{
		Now:           func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) },
		NewID:         func() string { return "user-1" },
		NewToken:      func() (string, error) { return "tok-1", nil },
		NewPassword:   func() (string, error) { return "generated-pass", nil },
		TTL:           time.Hour,
		AdminUsername: "graefik",
	}
}

func TestAuthService_EnsureAdmin_CreeSiVide(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	users.EXPECT().Count(mock.Anything).Return(0, nil).Once()
	hasher.EXPECT().Hash("generated-pass").Return("hashed", nil).Once()
	users.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == "graefik" && u.PasswordHash == "hashed" && u.ID == "user-1"
	})).Return(nil).Once()

	svc := service.NewAuthService(users, portmocks.NewMockSessionRepository(t), hasher, authCfg())

	res, err := svc.EnsureAdmin(context.Background())

	require.NoError(t, err)
	assert.True(t, res.Created)
	assert.Equal(t, "graefik", res.Username)
	assert.Equal(t, "generated-pass", res.GeneratedPassword)
}

func TestAuthService_EnsureAdmin_IgnoreSiExiste(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	users.EXPECT().Count(mock.Anything).Return(1, nil).Once()

	svc := service.NewAuthService(users, portmocks.NewMockSessionRepository(t), portmocks.NewMockPasswordHasher(t), authCfg())

	res, err := svc.EnsureAdmin(context.Background())

	require.NoError(t, err)
	assert.False(t, res.Created)
	assert.Empty(t, res.GeneratedPassword)
}

func TestAuthService_Login_OK(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	sessions := portmocks.NewMockSessionRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	user := &domain.User{ID: "user-1", Username: "graefik", PasswordHash: "hashed"}
	users.EXPECT().FindByUsername(mock.Anything, "graefik").Return(user, nil).Once()
	hasher.EXPECT().Compare("hashed", "secret").Return(nil).Once()
	sessions.EXPECT().Create(mock.Anything, mock.MatchedBy(func(s *domain.Session) bool {
		return s.Token == "tok-1" && s.UserID == "user-1"
	})).Return(nil).Once()

	svc := service.NewAuthService(users, sessions, hasher, authCfg())

	session, err := svc.Login(context.Background(), "graefik", "secret")

	require.NoError(t, err)
	assert.Equal(t, "tok-1", session.Token)
	assert.Equal(t, "user-1", session.UserID)
}

func TestAuthService_Login_UtilisateurInconnu(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	users.EXPECT().FindByUsername(mock.Anything, "ghost").Return(nil, domain.ErrUserNotFound).Once()

	svc := service.NewAuthService(users, portmocks.NewMockSessionRepository(t), portmocks.NewMockPasswordHasher(t), authCfg())

	_, err := svc.Login(context.Background(), "ghost", "x")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestAuthService_Login_MauvaisMotDePasse(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	user := &domain.User{ID: "user-1", Username: "graefik", PasswordHash: "hashed"}
	users.EXPECT().FindByUsername(mock.Anything, "graefik").Return(user, nil).Once()
	hasher.EXPECT().Compare("hashed", "bad").Return(errors.New("mismatch")).Once()

	svc := service.NewAuthService(users, portmocks.NewMockSessionRepository(t), hasher, authCfg())

	_, err := svc.Login(context.Background(), "graefik", "bad")

	assert.ErrorIs(t, err, domain.ErrInvalidCredentials)
}

func TestAuthService_Authenticate_Valide(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	sessions := portmocks.NewMockSessionRepository(t)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	session := &domain.Session{Token: "tok-1", UserID: "user-1", ExpiresAt: now.Add(time.Hour)}
	sessions.EXPECT().FindByToken(mock.Anything, "tok-1").Return(session, nil).Once()
	sessions.EXPECT().UpdateExpiry(mock.Anything, "tok-1", mock.Anything).Return(nil).Once()
	users.EXPECT().FindByID(mock.Anything, "user-1").Return(&domain.User{ID: "user-1", Username: "graefik"}, nil).Once()

	svc := service.NewAuthService(users, sessions, portmocks.NewMockPasswordHasher(t), authCfg())

	user, err := svc.Authenticate(context.Background(), "tok-1")

	require.NoError(t, err)
	assert.Equal(t, "graefik", user.Username)
}

func TestAuthService_Authenticate_Expiree(t *testing.T) {
	sessions := portmocks.NewMockSessionRepository(t)
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	expired := &domain.Session{Token: "tok-1", UserID: "user-1", ExpiresAt: now.Add(-time.Minute)}
	sessions.EXPECT().FindByToken(mock.Anything, "tok-1").Return(expired, nil).Once()
	sessions.EXPECT().Delete(mock.Anything, "tok-1").Return(nil).Once()

	svc := service.NewAuthService(portmocks.NewMockUserRepository(t), sessions, portmocks.NewMockPasswordHasher(t), authCfg())

	_, err := svc.Authenticate(context.Background(), "tok-1")

	assert.ErrorIs(t, err, domain.ErrSessionInvalid)
}

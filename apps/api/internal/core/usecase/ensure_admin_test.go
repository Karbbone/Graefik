package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/usecase"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
)

func baseCfg(users *portmocks.MockUserRepository, hasher *portmocks.MockPasswordHasher) usecase.EnsureAdminConfig {
	return usecase.EnsureAdminConfig{
		Users:         users,
		Hasher:        hasher,
		Now:           fixedNow,
		NewID:         func() string { return "user-1" },
		NewPassword:   func() (string, error) { return "generated-pass", nil },
		AdminUsername: "graefik",
	}
}

func TestEnsureAdmin_ProdCreation(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	users.EXPECT().Count(mock.Anything).Return(0, nil).Once()
	hasher.EXPECT().Hash("generated-pass").Return("hashed", nil).Once()
	users.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == "graefik" && u.PasswordHash == "hashed"
	})).Return(nil).Once()

	uc := usecase.NewEnsureAdmin(baseCfg(users, hasher))

	res, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.True(t, res.Created)
	assert.Equal(t, "generated-pass", res.GeneratedPassword)
}

func TestEnsureAdmin_ProdExistant(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	users.EXPECT().Count(mock.Anything).Return(1, nil).Once()

	uc := usecase.NewEnsureAdmin(baseCfg(users, portmocks.NewMockPasswordHasher(t)))

	res, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.False(t, res.Created)
	assert.Empty(t, res.GeneratedPassword)
}

func TestEnsureAdmin_DevCreation(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	cfg := baseCfg(users, hasher)
	cfg.DevMode = true
	cfg.DevPassword = "dev-pass"

	users.EXPECT().Count(mock.Anything).Return(0, nil).Once()
	hasher.EXPECT().Hash("dev-pass").Return("hashed", nil).Once()
	users.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()

	uc := usecase.NewEnsureAdmin(cfg)

	res, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.True(t, res.Created)
	assert.Equal(t, "dev-pass", res.GeneratedPassword)
}

func TestEnsureAdmin_DevReinitialise(t *testing.T) {
	users := portmocks.NewMockUserRepository(t)
	hasher := portmocks.NewMockPasswordHasher(t)

	cfg := baseCfg(users, hasher)
	cfg.DevMode = true
	cfg.DevPassword = "dev-pass"

	existing := &domain.User{ID: "user-1", Username: "graefik", PasswordHash: "old"}
	users.EXPECT().Count(mock.Anything).Return(1, nil).Once()
	users.EXPECT().FindByUsername(mock.Anything, "graefik").Return(existing, nil).Once()
	hasher.EXPECT().Hash("dev-pass").Return("newhash", nil).Once()
	users.EXPECT().UpdatePassword(mock.Anything, "user-1", "newhash").Return(nil).Once()

	uc := usecase.NewEnsureAdmin(cfg)

	res, err := uc.Execute(context.Background())

	require.NoError(t, err)
	assert.False(t, res.Created)
	assert.Equal(t, "dev-pass", res.GeneratedPassword)
}

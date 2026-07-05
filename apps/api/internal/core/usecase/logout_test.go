package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Karbbone/Graefik/apps/api/internal/core/usecase"
	portmocks "github.com/Karbbone/Graefik/apps/api/internal/mocks/port"
)

func TestLogout(t *testing.T) {
	sessions := portmocks.NewMockSessionRepository(t)
	sessions.EXPECT().Delete(mock.Anything, "tok-1").Return(nil).Once()

	uc := usecase.NewLogout(sessions)

	require.NoError(t, uc.Execute(context.Background(), "tok-1"))
}

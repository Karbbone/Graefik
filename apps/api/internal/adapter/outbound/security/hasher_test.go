package security_test

import (
	"testing"

	"github.com/Karbbone/Graefik/apps/api/internal/adapter/outbound/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBcryptHasher(t *testing.T) {
	h := security.NewBcryptHasher()

	hash, err := h.Hash("mon-mot-de-passe")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, "mon-mot-de-passe", hash)

	assert.NoError(t, h.Compare(hash, "mon-mot-de-passe"))
	assert.Error(t, h.Compare(hash, "mauvais"))
}

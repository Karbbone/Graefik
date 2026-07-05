// Package security est un adaptateur outbound : hachage de mots de passe (bcrypt).
package security

import (
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
	"golang.org/x/crypto/bcrypt"
)

// BcryptHasher implémente port.PasswordHasher via bcrypt.
type BcryptHasher struct {
	cost int
}

var _ port.PasswordHasher = (*BcryptHasher)(nil)

// NewBcryptHasher crée un hasher au coût par défaut de bcrypt.
func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{cost: bcrypt.DefaultCost}
}

// Hash produit le hash bcrypt du mot de passe (sel intégré).
func (h *BcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare vérifie le mot de passe contre le hash (temps constant).
func (h *BcryptHasher) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

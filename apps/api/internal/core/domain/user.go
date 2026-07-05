// Package domain contient le cœur métier : entités, value objects, règles et
// erreurs. Il ne dépend d'AUCUN framework ni détail technique (ni Echo, ni BDD).
package domain

import (
	"errors"
	"time"
)

// ErrUserNotFound est renvoyée par le repository quand aucun utilisateur ne correspond.
var ErrUserNotFound = errors.New("utilisateur introuvable")

// User est l'entité représentant un compte.
type User struct {
	ID           string
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

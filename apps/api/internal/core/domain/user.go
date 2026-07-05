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

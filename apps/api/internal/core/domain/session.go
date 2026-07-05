package domain

import (
	"errors"
	"time"
)

// Erreurs d'authentification.
var (
	// ErrInvalidCredentials : identifiant ou mot de passe incorrect.
	// Volontairement générique pour ne pas divulguer quel champ est faux.
	ErrInvalidCredentials = errors.New("identifiants invalides")
	// ErrSessionInvalid : session absente, inconnue ou expirée.
	ErrSessionInvalid = errors.New("session invalide ou expirée")
)

// Session représente une session authentifiée persistée.
type Session struct {
	Token     string
	UserID    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// IsExpired indique si la session est expirée à l'instant donné.
func (s *Session) IsExpired(now time.Time) bool {
	return !now.Before(s.ExpiresAt)
}

package service

import (
	"context"
	"errors"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// AuthConfig regroupe les dépendances injectées du service d'authentification.
type AuthConfig struct {
	Now           func() time.Time
	NewID         func() string             // identifiants d'utilisateur
	NewToken      func() (string, error)    // tokens de session (crypto-aléatoires)
	NewPassword   func() (string, error)    // mot de passe admin généré
	TTL           time.Duration             // durée de vie d'une session
	AdminUsername string
}

// AuthService implémente le port inbound port.AuthService.
type AuthService struct {
	users    port.UserRepository
	sessions port.SessionRepository
	hasher   port.PasswordHasher
	cfg      AuthConfig
}

var _ port.AuthService = (*AuthService)(nil)

// NewAuthService construit le service avec ses ports outbound et sa configuration.
func NewAuthService(users port.UserRepository, sessions port.SessionRepository, hasher port.PasswordHasher, cfg AuthConfig) *AuthService {
	return &AuthService{users: users, sessions: sessions, hasher: hasher, cfg: cfg}
}

// EnsureAdmin crée l'admin initial si aucun compte n'existe.
func (s *AuthService) EnsureAdmin(ctx context.Context) (port.BootstrapResult, error) {
	count, err := s.users.Count(ctx)
	if err != nil {
		return port.BootstrapResult{}, err
	}
	if count > 0 {
		return port.BootstrapResult{Created: false}, nil
	}

	password, err := s.cfg.NewPassword()
	if err != nil {
		return port.BootstrapResult{}, err
	}
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return port.BootstrapResult{}, err
	}

	user := &domain.User{
		ID:           s.cfg.NewID(),
		Username:     s.cfg.AdminUsername,
		PasswordHash: hash,
		CreatedAt:    s.cfg.Now(),
	}
	if err := s.users.Create(ctx, user); err != nil {
		return port.BootstrapResult{}, err
	}

	return port.BootstrapResult{
		Created:           true,
		Username:          user.Username,
		GeneratedPassword: password,
	}, nil
}

// Login vérifie les identifiants et crée une session.
func (s *AuthService) Login(ctx context.Context, username, password string) (*domain.Session, error) {
	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			// On hache tout de même une valeur bidon serait idéal contre le timing ;
			// ici on renvoie une erreur générique.
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := s.cfg.NewToken()
	if err != nil {
		return nil, err
	}
	now := s.cfg.Now()
	session := &domain.Session{
		Token:     token,
		UserID:    user.ID,
		CreatedAt: now,
		ExpiresAt: now.Add(s.cfg.TTL),
	}
	if err := s.sessions.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

// Logout invalide la session (idempotent).
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessions.Delete(ctx, token)
}

// Authenticate valide un token, prolonge la session (glissante) et renvoie l'utilisateur.
func (s *AuthService) Authenticate(ctx context.Context, token string) (*domain.User, error) {
	session, err := s.sessions.FindByToken(ctx, token)
	if err != nil {
		return nil, domain.ErrSessionInvalid
	}

	now := s.cfg.Now()
	if session.IsExpired(now) {
		_ = s.sessions.Delete(ctx, token)
		return nil, domain.ErrSessionInvalid
	}

	// Session glissante : on repousse l'expiration (best-effort).
	newExpiry := now.Add(s.cfg.TTL)
	_ = s.sessions.UpdateExpiry(ctx, token, newExpiry)

	user, err := s.users.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, domain.ErrSessionInvalid
	}
	return user, nil
}

package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
)

// BootstrapResult décrit le résultat de l'amorçage de l'admin.
// GeneratedPassword n'est renseigné que lorsqu'il y a un mot de passe à afficher.
type BootstrapResult struct {
	Created           bool
	Username          string
	GeneratedPassword string
}

// EnsureAdminConfig regroupe les dépendances du cas d'usage EnsureAdmin.
type EnsureAdminConfig struct {
	Users         port.UserRepository
	Hasher        port.PasswordHasher
	Now           func() time.Time
	NewID         func() string
	NewPassword   func() (string, error)
	AdminUsername string
	// DevMode : en dev, le mot de passe admin est fixé à DevPassword et
	// réinitialisé/réaffiché à chaque démarrage. En prod, généré une fois.
	DevMode     bool
	DevPassword string
}

// EnsureAdmin garantit l'existence de l'admin (appelé au démarrage).
type EnsureAdmin struct {
	cfg EnsureAdminConfig
}

// NewEnsureAdmin construit le cas d'usage EnsureAdmin.
func NewEnsureAdmin(cfg EnsureAdminConfig) *EnsureAdmin {
	return &EnsureAdmin{cfg: cfg}
}

// Execute :
//   - aucun compte : création (dev = DevPassword, prod = aléatoire) ;
//   - compte existant en dev : réinitialise le mot de passe à DevPassword ;
//   - compte existant en prod : aucune action.
func (uc *EnsureAdmin) Execute(ctx context.Context) (BootstrapResult, error) {
	count, err := uc.cfg.Users.Count(ctx)
	if err != nil {
		return BootstrapResult{}, err
	}

	if count == 0 {
		return uc.createAdmin(ctx)
	}
	if !uc.cfg.DevMode {
		return BootstrapResult{Created: false}, nil
	}
	return uc.resetAdminPassword(ctx)
}

func (uc *EnsureAdmin) createAdmin(ctx context.Context) (BootstrapResult, error) {
	password := uc.cfg.DevPassword
	if !uc.cfg.DevMode {
		var err error
		if password, err = uc.cfg.NewPassword(); err != nil {
			return BootstrapResult{}, err
		}
	}

	hash, err := uc.cfg.Hasher.Hash(password)
	if err != nil {
		return BootstrapResult{}, err
	}

	user := &domain.User{
		ID:           uc.cfg.NewID(),
		Username:     uc.cfg.AdminUsername,
		PasswordHash: hash,
		CreatedAt:    uc.cfg.Now(),
	}
	if err := uc.cfg.Users.Create(ctx, user); err != nil {
		return BootstrapResult{}, err
	}

	return BootstrapResult{Created: true, Username: user.Username, GeneratedPassword: password}, nil
}

func (uc *EnsureAdmin) resetAdminPassword(ctx context.Context) (BootstrapResult, error) {
	user, err := uc.cfg.Users.FindByUsername(ctx, uc.cfg.AdminUsername)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return BootstrapResult{Created: false}, nil
		}
		return BootstrapResult{}, err
	}

	hash, err := uc.cfg.Hasher.Hash(uc.cfg.DevPassword)
	if err != nil {
		return BootstrapResult{}, err
	}
	if err := uc.cfg.Users.UpdatePassword(ctx, user.ID, hash); err != nil {
		return BootstrapResult{}, err
	}

	return BootstrapResult{Created: false, Username: user.Username, GeneratedPassword: uc.cfg.DevPassword}, nil
}

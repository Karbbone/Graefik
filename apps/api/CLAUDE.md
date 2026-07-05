# CLAUDE.md — Backend (`apps/api`)

Backend **Go 1.26 + Echo v5** en **architecture hexagonale** (ports & adapters). Lis la [racine CLAUDE.md](../../CLAUDE.md) pour les commandes Docker.

## Principe : la dépendance pointe vers l'intérieur

```
        adapter (HTTP, mémoire…)
             │  dépend de
             ▼
      port (interfaces)
             │
             ▼
        domain (cœur)
```

- Le **domain** ne dépend de RIEN (ni Echo, ni BDD, ni autre package interne).
- Les **services** (use cases) dépendent du domain et des **ports**, jamais des adapters.
- Les **adapters** dépendent des ports. Ils sont interchangeables (mémoire → Postgres) sans toucher au cœur.
- Le câblage des implémentations concrètes se fait UNIQUEMENT dans `cmd/api/main.go` (composition root).

## Structure

```
apps/api/
├─ cmd/api/main.go                 # composition root : instancie et câble tout
└─ internal/
   ├─ core/                        # L'HEXAGONE — aucun import de framework
   │  ├─ domain/                   # entités, value objects, erreurs & règles métier
   │  ├─ port/                     # interfaces : inbound.go (use cases), outbound.go (repos/gateways)
   │  └─ service/                  # use cases : implémentent les ports inbound
   ├─ adapter/
   │  ├─ inbound/http/             # driving : handlers Echo, DTO, routes (router.go)
   │  └─ outbound/memory/          # driven : impl. en mémoire des ports outbound
   ├─ platform/                    # config, (futur) logger, pool DB
   └─ mocks/port/                  # mocks générés (mockery) — NE PAS éditer à la main
```

## Règles d'import (à respecter absolument)

| Package | Peut importer | Ne doit JAMAIS importer |
|---------|---------------|-------------------------|
| `core/domain` | (rien du projet) | port, service, adapter, echo |
| `core/port` | `core/domain` | service, adapter, echo |
| `core/service` | `core/domain`, `core/port` | adapter, echo |
| `adapter/*` | `core/port`, `core/domain`, echo | un autre adapter concret |
| `cmd/api` | tout (câblage) | — |

## Où ajouter du code

- **Nouvelle règle métier** → `core/domain` (+ test). Ex. une méthode sur `Task`, une validation.
- **Nouveau cas d'usage** → ajoute la méthode à l'interface dans `core/port/inbound.go`, implémente-la dans `core/service`, expose-la via un handler dans `adapter/inbound/http`, câble dans `cmd/api/main.go`.
- **Nouvelle dépendance externe** (BDD, API tierce…) → déclare un port dans `core/port/outbound.go`, écris l'adaptateur dans `adapter/outbound/<techno>/`, injecte-le dans `main.go`. Le cœur ne change pas.
- **Nouvel endpoint HTTP** → handler dans `adapter/inbound/http`, route dans `router.go`.

## Tests

Stack : **testify** (`assert`/`require`) + **mockery** (mocks des ports). C'est le combo standard de l'écosystème Go.

- Tests **table-driven** quand c'est pertinent (voir `core/domain/task_test.go`).
- `core/domain` → tests unitaires purs.
- `core/service` → tests unitaires avec le repository **mocké** (`internal/mocks/port`) : on teste les cas OK, les erreurs métier, les erreurs du repo.
- `adapter/outbound/memory` → tests de l'implémentation réelle.
- `adapter/inbound/http` → tests via `e.ServeHTTP(rec, req)` + `httptest`, service mocké.

Lancer les tests :

```bash
docker compose exec -T api sh -c "cd /app && go test ./..."
```

### Régénérer les mocks

Les mocks vivent dans `internal/mocks/port/` et sont **commités**. Après avoir modifié une interface dans `core/port`, régénère-les :

```bash
docker compose exec -T api sh -c "cd /app && go run github.com/vektra/mockery/v3@v3.7.1"
```

Config dans `.mockery.yaml` (ajoute-y toute nouvelle interface à mocker).

## Authentification

Implémentée en respectant l'hexagonal :

- **Domaine** : `domain/user.go` (`User`), `domain/session.go` (`Session`, `ErrInvalidCredentials`, `ErrSessionInvalid`).
- **Ports outbound** (`core/port/outbound.go`) : `UserRepository`, `SessionRepository`, `PasswordHasher`.
- **Port inbound** (`core/port/inbound.go`) : `AuthService` (`EnsureAdmin`, `Login`, `Logout`, `Authenticate`).
- **Use case** : `core/service/auth_service.go` (dépendances injectées via `AuthConfig` : horloge, générateurs d'ID/token/mot de passe, TTL).
- **Adaptateurs outbound** : `adapter/outbound/sqlite/` (`db.go` + migrations, `user_repository.go`, `session_repository.go`, driver pur-Go `modernc.org/sqlite`) et `adapter/outbound/security/hasher.go` (bcrypt).
- **Adaptateur inbound** : `adapter/inbound/http/auth_handler.go`, `middleware.go` (`RequireAuth`), `cookie.go`, `login_limiter.go` (verrou anti-brute-force en mémoire, 5 échecs → 1 min).
- **Génération de secrets** : `platform/secret.go` (`GenerateToken`, `GeneratePassword`).
- **Bootstrap** : `cmd/api/main.go` appelle `EnsureAdmin` au démarrage.
  - **Prod** : crée l'admin avec un mot de passe aléatoire, affiché **une fois**.
  - **Dev** (`DevMode`, via `APP_ENV=development`) : réinitialise le mot de passe admin à `DEV_ADMIN_PASSWORD` (défaut `graefik`) et le **réaffiche à chaque démarrage** (`UpdatePassword`). `EnsureAdmin` renvoie `GeneratedPassword` non vide dès qu'il y a un mot de passe à afficher.

**Config** (env, voir `platform/config.go`) : `DB_PATH` (défaut `/data/graefik.db`), `COOKIE_NAME`, `COOKIE_SECURE` (`auto`/`true`/`false`), `SESSION_TTL`.

**Cookie** : `graefik_session`, HttpOnly, SameSite=Lax, Secure auto, 7 j glissants.

> Les mocks (`internal/mocks/port`) couvrent aussi `UserRepository`, `SessionRepository`, `PasswordHasher`, `AuthService` — régénère-les après toute modif d'interface (voir plus haut).

## Lint & format

Outil unique : **golangci-lint v2** (installé dans l'image dev, config `.golangci.yml`).

```bash
docker compose exec -T api sh -c "cd /app && golangci-lint fmt"   # formate (gofumpt + goimports)
docker compose exec -T api sh -c "cd /app && golangci-lint run"   # lint (0 issue exigé)
```

- **Formatage** : gofumpt (plus strict que gofmt) + goimports (imports locaux `github.com/Karbbone/Graefik` regroupés).
- **Linters** : jeu standard (govet, staticcheck, errcheck, ineffassign, unused) + revive, gocritic, errorlint, unconvert, bodyclose. Les mocks générés et les `_test.go` sont exclus/assouplis.
- Règles apprises : pas de redéfinition de builtins (`max`, `min`…), commentaires de doc sur tout l'exporté, pas d'`os.Exit` après un `defer` (logique dans `run() error`). `misspell` est désactivé (commentaires en français).

## Conventions Go

- Handlers Echo v5 : signature `func(c *echo.Context) error` (pointeur, ≠ v4).
- Erreurs métier : `var Err… = errors.New(...)` dans `domain`, comparées avec `errors.Is`.
- Injection de dépendances par constructeur (`New…`). Horloge et génération d'ID injectées comme fonctions (`func() time.Time`, `func() string`) pour la testabilité.
- Assertion de conformité aux ports : `var _ port.X = (*Impl)(nil)`.
- Hot reload assuré par **Air** (`.air.toml`, build de `./cmd/api`).

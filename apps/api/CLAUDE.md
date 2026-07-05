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

## Conventions Go

- Handlers Echo v5 : signature `func(c *echo.Context) error` (pointeur, ≠ v4).
- Erreurs métier : `var Err… = errors.New(...)` dans `domain`, comparées avec `errors.Is`.
- Injection de dépendances par constructeur (`New…`). Horloge et génération d'ID injectées comme fonctions (`func() time.Time`, `func() string`) pour la testabilité.
- Assertion de conformité aux ports : `var _ port.X = (*Impl)(nil)`.
- Hot reload assuré par **Air** (`.air.toml`, build de `./cmd/api`).

# CLAUDE.md — Graefik

Guide destiné à Claude Code pour travailler dans ce dépôt. Lis-le en entier avant toute modification.

## Contexte produit — que fait Graefik ?

**Graefik** est une application de **monitoring et de data-visualisation self-hosted**, spécialisée dans les **reverse proxies** (Traefik en premier). Le nom est un mot-valise : **Graf**ana + Tra**efik**, et un clin d'œil à « graphique ». Prononcé « graphique ».

### Le problème
Un reverse proxy comme Traefik expose énormément de métriques (routeurs, services, volume de requêtes, latences, codes HTTP, erreurs…), mais les exploiter demande aujourd'hui d'assembler et maintenir soi-même une stack lourde (Prometheus + Grafana + configuration). C'est disproportionné pour « juste » voir l'état de ses routeurs.

### La proposition
Graefik = **« un Grafana en un »**, prêt à l'emploi et spécialisé reverse proxy. On le déploie comme **un simple conteneur**, à côté de son Traefik (exactement dans le même esprit self-hosted que Traefik lui-même), et on obtient directement des tableaux de bord des différents routeurs/services.

### Ce que fait (et fera) l'application
1. **Se connecter** à un reverse proxy (Traefik d'abord) et **récupérer ses métriques**.
2. **Stocker** ces métriques dans le temps — base de données **time-series** (candidat : InfluxDB / « flux », ou autre TSDB — *décision à venir, non figée*).
3. **Visualiser** les statistiques par **routeur** et par **service** : trafic, latence, taux d'erreur, codes HTTP, etc.
4. **Alerter** (à terme) : seuils, notifications.

### Cible de déploiement
**Self-hosted** chez l'utilisateur final, packagé en conteneur (image Docker / `docker compose`), simple à lancer et à configurer — la facilité de déploiement est un objectif produit central.

### Vocabulaire du domaine (ubiquitous language)
| Terme | Sens dans Graefik |
|-------|-------------------|
| **Source** (data source) | Un reverse proxy surveillé (ex. une instance Traefik) d'où proviennent les métriques |
| **Router** | Un routeur du reverse proxy (règle d'entrée → service) dont on affiche les stats |
| **Service** | La cible d'un routeur (backend) |
| **Métrique** | Une mesure horodatée (req/s, latence, code HTTP, erreurs…) |
| **Dashboard** | Une vue agrégeant les métriques d'un ou plusieurs routeurs/services |
| **Alerte** | Une règle déclenchant une notification quand une métrique franchit un seuil |

### Comment la vision se traduit dans l'architecture
L'hexagonal (voir [apps/api/CLAUDE.md](apps/api/CLAUDE.md)) est choisi précisément pour ce contexte :
- La **récupération des métriques** sera un **port outbound** (ex. `MetricsSource`) avec un adaptateur Traefik — d'autres proxies pourront suivre sans toucher au cœur.
- Le **stockage time-series** sera un autre **port outbound** (ex. `MetricsStore`) avec un adaptateur InfluxDB/TSDB — la techno de BDD reste ainsi remplaçable.
- Le domaine (routeurs, métriques, dashboards, alertes) reste pur et indépendant du proxy comme de la BDD.

> ⚠️ **État actuel** : seule l'**authentification** est implémentée (voir plus bas). Les features métier (sources, métriques, dashboards, alertes) restent à construire.

## Vue d'ensemble

Graefik est un **monorepo** :

- `apps/api` — backend **Go 1.26 + Echo v5**, architecture **hexagonale** (ports & adapters). Voir [apps/api/CLAUDE.md](apps/api/CLAUDE.md).
- `apps/web` — frontend **React 19 + Vite 8** (SPA TypeScript), architecture **feature-based**. Voir [apps/web/CLAUDE.md](apps/web/CLAUDE.md).
- `packages/` — paquets partagés (vide pour l'instant).

Orchestration : **Turborepo** + workspaces **bun**. Développement 100 % conteneurisé (**Docker**) avec hot reload : rien à installer sur la machine hôte à part Docker.

## Règle d'or : tout passe par Docker

**Go et bun ne sont PAS installés sur la machine hôte.** N'exécute jamais `go`, `bun`, `npm` ou `node` directement sur l'hôte : ça échouera. Utilise les conteneurs.

```bash
docker compose up --build      # lance toute la stack (web + api) avec hot reload
docker compose down            # arrête tout
docker compose logs -f api     # logs du backend (Air)
docker compose logs -f web     # logs du frontend (Vite)
```

Exécuter une commande de toolchain (tests, lint, build, génération) = `docker compose exec` dans le bon service :

```bash
# Backend (service "api", contient Go)
docker compose exec -T api sh -c "cd /app && go test ./..."
docker compose exec -T api sh -c "cd /app && go vet ./..."

# Frontend (service "web", contient bun)
docker compose exec -T web sh -c "cd /app && bun run test"
docker compose exec -T web sh -c "cd /app && bun run lint"
docker compose exec -T web sh -c "cd /app && bun run check-types"
docker compose exec -T web sh -c "cd /app && bun run format"
```

### Qualité de code (formatage + lint)

| | Frontend (`apps/web`) | Backend (`apps/api`) |
|--|----------------------|----------------------|
| Formateur | **Biome** (`bun run format`) | **gofumpt** via `golangci-lint fmt` |
| Linter | **Biome** + **ESLint** (boundaries) → `bun run lint` | **golangci-lint** → `golangci-lint run` |

- Front : `bun run format` corrige (Biome), `bun run lint` vérifie (Biome + la règle d'architecture ESLint). Style imposé : guillemets doubles, points-virgules, imports triés.
- Back : `golangci-lint fmt` formate, `golangci-lint run` lint (govet/staticcheck/errcheck/revive/gocritic/errorlint…). Config `apps/api/.golangci.yml`.
- Front : **Knip** (`bun run knip`) détecte fichiers/dépendances/exports morts.

### Intégration continue

`.github/workflows/ci.yml` rejoue tout à chaque push / PR sur `main`, en deux jobs (les runners installent Go et bun directement, **pas de Docker en CI**) :

- **Backend** : mocks à jour (mockery + `git diff`), `golangci-lint fmt --diff`, `golangci-lint run`, `go test`, `go build`.
- **Frontend** : `bun run lint` (Biome + ESLint), `check-types`, `test` (Vitest), `knip`, `build`.

> Si un conteneur n'est pas démarré, lance `docker compose up -d` d'abord.

## URLs (dev)

| Service | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| API | http://localhost:${API_PORT}/api/health |

⚠️ **Ports** : le port hôte de l'API est configurable via `.env` (`API_PORT`, `WEB_PORT`). Sur cette machine, `8080` est déjà pris par un autre projet, donc `.env` mappe l'API sur **8081** (le port *interne* du conteneur reste 8080, le proxy Vite est inchangé). Ne pas committer `.env` (il est gitignoré) ; `.env.example` documente les variables.

## Authentification (self-hosted)

Modèle « à la Jenkins » — comportement différent selon l'environnement :

- **Prod** (`APP_ENV` ≠ `development`) : au **premier lancement**, création de l'admin **`graefik`** avec un mot de passe **aléatoire**, affiché **une seule fois** dans les logs.
- **Dev** (`APP_ENV=development`, cas du `docker compose`) : le mot de passe admin est **fixé** (`DEV_ADMIN_PASSWORD`, défaut `graefik`), **réinitialisé et réaffiché à chaque démarrage** — pratique pour se reconnecter sans fouiller.

```bash
docker compose logs api | grep -A6 "administrateur"
```

- Session par **cookie HttpOnly** (`graefik_session`), stockée en **SQLite** (persistée sur le volume `graefik-data`), expiration **7 jours glissants**, flag `Secure` **auto** (HTTPS/`X-Forwarded-Proto`, override `COOKIE_SECURE`).
- Endpoints : `POST /api/auth/login`, `POST /api/auth/logout`, `GET /api/auth/me`. `/api/health` est public ; les futures routes métier seront protégées par session (middleware `RequireAuth`).
- Détails backend : [apps/api/CLAUDE.md](apps/api/CLAUDE.md) · flux frontend : [apps/web/CLAUDE.md](apps/web/CLAUDE.md).

## Turborepo

`turbo.json` définit les tâches `dev`, `build`, `lint`, `check-types`, `test`, `clean`. Chaque app expose ces scripts dans son `package.json`. Depuis la racine (dans un environnement qui aurait bun) : `bun run build`, `bun run lint`, etc. En pratique on passe par Docker (voir ci-dessus).

## Conventions

- **Langue** : commentaires, messages de commit et documentation en **français** (avec accents corrects). Le code (identifiants) reste en anglais.
- **Commits** : format court type `type: description` (`feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`). Ne committer/pusher que sur demande explicite.
- **Tests** : toute nouvelle logique métier vient avec ses tests. Le back et le front ont chacun leur stratégie (voir leurs CLAUDE.md respectifs).
- **Ne jamais** committer `.env`, `node_modules`, `dist`, `tmp/`, `bin/`.

## Où travailler

| Tu veux… | Va dans… |
|----------|----------|
| Ajouter une règle métier / un endpoint backend | `apps/api` → lis [apps/api/CLAUDE.md](apps/api/CLAUDE.md) |
| Ajouter un écran / une feature frontend | `apps/web` → lis [apps/web/CLAUDE.md](apps/web/CLAUDE.md) |
| Modifier l'orchestration dev | `docker-compose.yml`, `turbo.json` |

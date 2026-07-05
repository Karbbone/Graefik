# CLAUDE.md — Graefik

Guide destiné à Claude Code pour travailler dans ce dépôt. Lis-le en entier avant toute modification.

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
```

> Si un conteneur n'est pas démarré, lance `docker compose up -d` d'abord.

## URLs (dev)

| Service | URL |
|---------|-----|
| Frontend | http://localhost:5173 |
| API | http://localhost:${API_PORT}/api/health |

⚠️ **Ports** : le port hôte de l'API est configurable via `.env` (`API_PORT`, `WEB_PORT`). Sur cette machine, `8080` est déjà pris par un autre projet, donc `.env` mappe l'API sur **8081** (le port *interne* du conteneur reste 8080, le proxy Vite est inchangé). Ne pas committer `.env` (il est gitignoré) ; `.env.example` documente les variables.

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

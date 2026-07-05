# Graefik

Monorepo **Go (Echo)** + **React (Vite)**, orchestré par **Turborepo** (gestionnaire de paquets **bun**), entièrement dockerisé avec **hot reload** — aucune installation de Go ni de Node/Bun requise sur ta machine, seul **Docker** est nécessaire.

## Stack

| Brique | Choix | Version |
|--------|-------|---------|
| Monorepo | Turborepo + bun workspaces | Turbo ^2.10 / bun 1.3 |
| Backend | Go + Echo | Go 1.26 / Echo v5 |
| Hot reload Go | Air | v1.65.2 |
| Frontend | React + Vite (SPA + TypeScript) | React 19 / Vite 8 |
| Conteneurs | Docker Compose | — |

## Structure

```
Graefik/
├─ apps/
│  ├─ api/                  # Backend Go (Echo v5)
│  │  ├─ main.go
│  │  ├─ internal/
│  │  │  ├─ server/         # Instance Echo, middlewares, routes
│  │  │  └─ handler/        # Handlers HTTP
│  │  ├─ .air.toml          # Config hot reload
│  │  ├─ Dockerfile.dev     # Image dev (Air)
│  │  └─ Dockerfile         # Image prod (multi-stage, distroless)
│  └─ web/                  # Frontend React (Vite)
│     ├─ src/
│     ├─ vite.config.ts     # Proxy /api -> backend
│     ├─ Dockerfile.dev     # Image dev (bun + Vite HMR)
│     └─ Dockerfile         # Image prod (build -> nginx)
├─ packages/                # Paquets partagés (à venir)
├─ docker-compose.yml       # Orchestration dev (les 2 services + hot reload)
├─ turbo.json               # Pipeline Turborepo
└─ package.json             # Workspaces bun + scripts
```

## Démarrage (recommandé : Docker, rien à installer)

```bash
docker compose up --build
# ou, raccourci :
bun run dev        # si bun est installé, sinon utilise la ligne ci-dessus
```

- Frontend : http://localhost:5173
- API : http://localhost:8080/api/health

Le hot reload est actif dans les deux sens :
- Modifie un fichier `.tsx` dans `apps/web` → Vite recharge le navigateur (HMR).
- Modifie un fichier `.go` dans `apps/api` → Air recompile et relance l'API.

Pour tout arrêter :

```bash
docker compose down
```

> Note Windows : le polling est activé (Vite `usePolling` et Air `poll`) pour que la
> détection des changements fonctionne à travers les bind mounts Docker.

## Développement natif (optionnel, si tu installes les toolchains)

Nécessite Go 1.26+, bun 1.3+ et Air installés localement.

```bash
bun install
bun run dev:native   # turbo run dev : lance api + web en parallèle
```

## Scripts Turborepo (racine)

| Commande | Effet |
|----------|-------|
| `bun run dev` | Lance la stack via Docker Compose |
| `bun run dev:native` | Lance api + web en natif via Turbo |
| `bun run build` | Build tous les apps |
| `bun run lint` | Lint tous les apps |
| `bun run check-types` | Vérification de types Go + TS |
| `bun run test` | Tests de tous les apps |

## Build de production

```bash
# Backend -> binaire dans une image distroless minimale
docker build -t graefik-api ./apps/api

# Frontend -> build statique servi par nginx (proxy /api inclus)
docker build -t graefik-web ./apps/web
```

## API

| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/api/health` | Contrôle de santé (statut, service, version Go, horodatage) |

Configuration via variables d'environnement (voir `.env.example`) :
- `PORT` (défaut `8080`)
- `APP_ENV`
- `CORS_ORIGINS` (liste séparée par des virgules, défaut `*`)

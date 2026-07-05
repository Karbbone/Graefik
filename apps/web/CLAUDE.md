# CLAUDE.md — Frontend (`apps/web`)

Frontend **React 19 + Vite 8** (SPA TypeScript) en **architecture feature-based**. Lis la [racine CLAUDE.md](../../CLAUDE.md) pour les commandes Docker.

## Principe : des features autonomes

Chaque **feature** est une tranche verticale isolée (son UI, ses appels réseau, ses hooks, ses types). Une feature :

- n'importe **jamais** l'intérieur d'une autre feature ;
- expose une **API publique** via son `index.ts` — c'est le seul point d'entrée pour l'extérieur ;
- peut importer du `shared/` (code transverse) et sa propre feature.

Ces règles sont **imposées automatiquement** par ESLint (`eslint-plugin-boundaries`). Une violation = erreur de lint.

## Structure

```
apps/web/src/
├─ app/                 # bootstrap : App racine, providers globaux, styles globaux
├─ pages/               # composent des features en écrans (branchés au routeur)
├─ features/            # 1 dossier autonome par feature
│  └─ <feature>/
│     ├─ api/           #   accès réseau (fonctions/hooks)
│     ├─ components/    #   UI spécifique à la feature
│     ├─ hooks/         #   hooks React de la feature
│     ├─ model/         #   types & logique métier
│     └─ index.ts       #   API PUBLIQUE (seuls ces exports sont importables ailleurs)
├─ shared/              # transverse réutilisable
│  ├─ lib/              #   utils, client HTTP (http.ts)
│  ├─ ui/               #   composants UI génériques (design system)
│  └─ hooks/, types/
├─ test/                # setup Vitest + handlers MSW
└─ main.tsx             # point d'entrée
```

## Règles de frontières (imposées par ESLint)

| Depuis | Peut importer |
|--------|---------------|
| `app` | `app`, `pages`, `feature` (via index), `shared` |
| `pages` | `pages`, `feature` (via index), `shared` |
| `feature` | `shared`, **sa propre** feature uniquement |
| `shared` | `shared` uniquement |

## Alias & imports

- Alias `@` → `src` (configuré dans `vite.config.ts` et `tsconfig.app.json`).
- Import d'une feature depuis l'extérieur : **toujours** via son index, ex. `import { TasksPanel } from '@/features/tasks'`. Jamais `@/features/tasks/components/...`.
- Import interne à une feature : chemins relatifs (`../hooks/useTasks`).

## Où ajouter du code

- **Nouvelle feature** → crée `src/features/<nom>/` avec `api/`, `components/`, `hooks/`, `model/`, et un `index.ts` qui exporte l'API publique. Compose-la dans une `page`.
- **Nouvel écran** → `src/pages/<Nom>Page.tsx`, qui assemble des features.
- **Code réutilisable transverse** → `src/shared/` (jamais dans une feature).
- **Provider global** (router, contexte, thème) → `src/app/`.

## Tests

Stack : **Vitest** + **@testing-library/react** + **MSW** (mock réseau).

- Setup global : `src/test/setup.ts` (jest-dom + cycle de vie MSW). Handlers par défaut : `src/test/msw/handlers.ts`.
- Les tests d'une feature vivent dans `src/features/<nom>/__tests__/`.
- On teste le **comportement** via les rôles/textes accessibles (`getByRole`, `findByText`), pas les détails d'implémentation.
- Pour un cas réseau spécifique, surcharge un handler dans le test : `server.use(http.get(...))`.

Commandes :

```bash
docker compose exec -T web sh -c "cd /app && bun run test"          # vitest run
docker compose exec -T web sh -c "cd /app && bun run coverage"      # couverture
docker compose exec -T web sh -c "cd /app && bun run lint"          # eslint + boundaries
docker compose exec -T web sh -c "cd /app && bun run check-types"   # tsc
```

## Conventions

- Composants fonctionnels + hooks. Un hook par préoccupation (`useTasks`, `useHealth`).
- Les appels réseau passent par `@/shared/lib/http` (`apiFetch`) et le proxy Vite `/api` → backend Go.
- Types explicites pour les données d'API (dossier `model/`).
- Accessibilité : `aria-label` sur les champs, textes de bouton clairs (facilite aussi les tests).
- Hot reload (HMR) assuré par Vite ; polling activé pour Docker/Windows.

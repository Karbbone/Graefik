#!/usr/bin/env bash
set -euo pipefail

ROOT="${CONTAINER_WORKSPACE_FOLDER:-/workspaces/Graefik}"
cd "$ROOT"

mkdir -p apps/api/tmp

LOG_FILE="/tmp/graefik-dev.log"

port_listening() {
  local port="$1"
  if command -v ss >/dev/null 2>&1; then
    ss -tlnH 2>/dev/null | grep -q ":${port} "
    return
  fi
  (echo >/dev/tcp/127.0.0.1/"$port") >/dev/null 2>&1
}

if port_listening 5173 || port_listening 8080; then
  echo "Stack dev déjà en cours (ports 5173 et/ou 8080 occupés)."
  exit 0
fi

# node_modules isolé dans un volume Linux : vérifie que les binaires sont présents.
if [[ ! -x node_modules/.bin/vite ]]; then
  echo "Dépendances manquantes — installation…"
  bun install
fi

# Variables dev container : priorité sur le .env hôte (souvent configuré pour docker-compose).
export APP_ENV="development"
export PORT="8080"
export CORS_ORIGINS="http://localhost:5173"
export DB_PATH="$ROOT/apps/api/tmp/graefik.db"
export TURBO_UI=stream

echo "Démarrage de l'API et du frontend…"
: >"$LOG_FILE"
nohup bun run dev:native >>"$LOG_FILE" 2>&1 &
disown

for _ in $(seq 1 60); do
  if port_listening 5173 && port_listening 8080; then
    echo ""
    echo "  Frontend : http://localhost:5173"
    echo "  API      : http://localhost:8080/api/health"
    echo "  Logs     : tail -f $LOG_FILE"
    echo ""
    exit 0
  fi

  if ! pgrep -f "turbo run dev" >/dev/null 2>&1 && ! pgrep -f "bun run dev:native" >/dev/null 2>&1; then
    echo "Échec du démarrage. Dernières lignes des logs :"
    echo ""
    tail -30 "$LOG_FILE"
    exit 1
  fi

  sleep 1
done

echo "Les serveurs mettent plus de temps que prévu à démarrer."
echo "Consulte les logs : tail -f $LOG_FILE"
tail -20 "$LOG_FILE"
exit 0

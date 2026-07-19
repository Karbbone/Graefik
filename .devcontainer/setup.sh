#!/usr/bin/env bash
set -euo pipefail

ROOT="${CONTAINER_WORKSPACE_FOLDER:-/workspaces/Graefik}"
cd "$ROOT"

mkdir -p apps/api/tmp

echo "Installation des dépendances JavaScript (Linux)…"
bun install

echo "Téléchargement des modules Go…"
go -C apps/api mod download

echo "Installation des outils Go (air, golangci-lint)…"
go install github.com/air-verse/air@v1.65.2 || true
curl -sSfL https://golangci-lint.run/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.12.2 || true

echo "Setup terminé."

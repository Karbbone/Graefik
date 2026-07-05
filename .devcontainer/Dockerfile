# Image de base des Dev Containers (Debian bookworm, utilisateur "vscode", git/curl).
# Go est ajouté par la feature devcontainer (voir devcontainer.json).
FROM mcr.microsoft.com/devcontainers/base:bookworm

# bun installé globalement (pour tous les utilisateurs).
ENV BUN_INSTALL=/usr/local/bun
ENV PATH=${BUN_INSTALL}/bin:${PATH}
RUN curl -fsSL https://bun.sh/install | bash -s "bun-v1.3.14"

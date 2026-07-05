# syntax=docker/dockerfile:1

# --- Etape de build ---
FROM oven/bun:1 AS build
WORKDIR /app

COPY package.json ./
RUN bun install

COPY . .
RUN bun run build

# --- Image finale : nginx sert les fichiers statiques ---
FROM nginx:1.27-alpine AS runtime
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

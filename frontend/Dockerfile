# Builder Node 8 desde GHCR (EOL pero disponible como mirror library)
FROM ghcr.io/library/node:8.17.0 AS build
WORKDIR /app

# Fix repos stretch para poder instalar toolchain (como ya tenías)
RUN set -eux; \
    sed -i -e 's/deb.debian.org/archive.debian.org/g' \
    -e 's|security.debian.org|archive.debian.org|g' \
    -e '/stretch-updates/d' /etc/apt/sources.list; \
    echo 'Acquire::Check-Valid-Until "false";' > /etc/apt/apt.conf.d/99no-check-valid-until; \
    apt-get update; \
    apt-get install -y --no-install-recommends python make g++ ca-certificates; \
    rm -rf /var/lib/apt/lists/*

COPY package*.json ./
RUN if [ -f package-lock.json ]; then npm ci --no-audit --no-fund; \
    else npm install --no-audit --no-fund; fi

COPY . .
RUN npm run build

# Runtime Nginx desde GHCR
FROM ghcr.io/library/nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80

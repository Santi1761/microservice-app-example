# Pipelines – Desarrollo (Build/Test/Push)

Cada microservicio tiene un workflow similar (`build-*.yml`):
- **Disparadores**: `push` a `develop`, `feature/**`, `release/**`, `main`; y `pull_request` a `develop`.
- **Pasos**: checkout → buildx → login Docker Hub → build & push.
- **Tagging**: `:latest` y `:<git-sha>`. (Opcional: publicar `:vX.Y.Z` en releases).

Servicios:
- `build-auth-api.yml` (context `./auth-api`)
- `build-users-api.yml` (context `./users-api`)
- `build-todos-api.yml` (context `./todos-api`)
- `build-frontend.yml` (context `./frontend`)
- `build-log-processor.yml` (context `./log-message-processor`)

Secretos usados:
- `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN` (a nivel de repo o environment).

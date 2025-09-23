# Pipelines – Despliegue (Staging/Producción)

## Staging
- Workflow: `deploy-staging.yml`
- Dispara en `push` a `develop` y `release/**` (o manual).
- Job con **environment: staging** → usa `SSH_HOST`, `SSH_USER`, `SSH_PRIVATE_KEY`.
- Acción remota: `docker compose -f ops/compose/docker-compose.stg.yml pull && up -d`.

## Producción
- Workflow: `deploy-prod.yml`
- Dispara en `push` a `main` (y/o `release: published`).
- **environment: production** con aprobación requerida.
- Acción remota: `docker compose -f ops/compose/docker-compose.prod.yml pull && up -d`.

> Las imágenes referenciadas en los `compose` usan `${DOCKERHUB_USER}` y `${TAG}` (por defecto `latest`).

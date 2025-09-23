# CI/CD

## Build (imágenes Docker)
- Workflows por microservicio en `.github/workflows/build-*.yml`.
- Etiquetado:
  - develop → dev-<shortSHA>
  - release/x.y.z → rc-x.y.z
  - main → x.y.z (release final)

## Deploy (por entorno)
- Ramas de operaciones: `env/dev`, `env/staging`, `env/prod`.
- Workflows: `deploy-dev.yml`, `deploy-staging.yml`, `deploy-prod.yml`.
- Cada deploy:
  1) Instala Ansible y community.docker
  2) Se conecta por SSH a la VM del entorno
  3) Renderiza y sube `.env` (con DOCKERHUB_USER + TAG según rama)
  4) Copia `docker-compose.yml` y `nginx.conf`
  5) `docker compose pull` + `up -d`

## Variables/Secrets por entorno
- SSH: `*_USER`, `*_HOST`, `*_SSH_KEY` (clave privada PEM)
- Azure (si se automatiza Terraform en CI): `ARM_*`
- Docker Hub: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`

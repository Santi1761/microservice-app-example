# Taller 1 — Automatización de Infraestructura con Pipelines
**Curso:** Ingeniería de Software V  
**Autores:** Santiago Arboleda Velasco, David Santiago Donneys

## 0. Resumen ejecutivo
Monorepo con microservicios (Frontend Vue, Users Spring Boot, Todos Node/Express, Auth Go, Log-Processor Python) orquestados con Docker Compose y desplegados en una VM Azure mediante GitHub Actions + Ansible. La infraestructura base (RG, red, SSH key, VM) se maneja con Terraform. Las imágenes se publican en Docker Hub. Se aplican dos patrones de nube: **Cache-Aside** (Redis) y **Circuit Breaker/Fallback** (Todos-API ante falla de Users-API).  
Se diferenciaron dos entornos: **staging** (Nginx :8080) y **producción** (Nginx :80).

---

## 1. Estrategia de branching (Desarrollo) — 2.5%
- Flujo **GitFlow simplificado**:
  - `feature/*` → integración en `develop`.
  - `develop` → `master` (release a producción).
- **Protecciones**:
  - `develop`/`master` con branch protection: PR obligatorio + status checks.
- **Evidencia**: captura de Branches y reglas en Settings → Branches.

## 2. Estrategia de branching (Operaciones/IaC) — 2.5%
- Ramas `infra/*` para cambios de Terraform (plan/apply por ambiente).
- Workflows dedicados: `infra-plan`, `infra-apply-staging`, `infra-apply-prod`.
- **Evidencia**: screenshots de los runs y del plan/apply.

---

## 3. Patrones de diseño de nube (mínimo 2) — 15%
### 3.1 Cache-Aside (Redis)
- Servicio: **todos-api** (Node).
- Mecanismo:
  1. Buscar primero en **Redis**.
  2. Si no está, consultar backend, **guardar** en Redis con TTL=60s y devolver.
- Configuración relevante:
  - `REDIS_URL=redis://redis:6379`
  - `CACHE_TTL_SECONDS=60`
- **Prueba**:
  1. `FLUSHALL` en `microservice-app-stg-redis-1`.
  2. GET `http://<IP>:8080/api/todos/` con token → aparecen claves en Redis (`KEYS *`, `TTL` ≈ 60).

### 3.2 Circuit Breaker / Fallback a caché
- Servicio: **todos-api**.
- Comportamiento:
  - Si `users-api` está caído, `todos-api` responde desde **caché** (si existe) o retorna error controlado (sin romper toda la app).
- **Prueba**:
  1. Calentar caché con una o más lecturas de `/api/todos/`.
  2. Parar `users-api` solo en STAGING.
  3. Repetir `/api/todos/`: se atiende desde caché o con error controlado.

---

## 4. Diagrama de arquitectura — 15%
- Archivo `docs/arquitectura.puml` (incluido en este repositorio; también disponible abajo).
- Componentes:
  - **VM Azure** (Ubuntu 22.04) con Docker/Compose.
  - **Stacks**: `microservice-app-stg` (8080) y `microservice-app` (80).
  - **Servicios**: nginx, frontend, users-api, todos-api, auth-api, redis, log-processor.
  - **CI/CD**: GitHub Actions (build, deploy, infra), Docker Hub, Terraform.
- Ruteo Nginx:
  - `/` → frontend
  - `/api/users` → users-api:8083
  - `/api/todos` → todos-api:8082
  - `/api/auth/*` y `/login` → auth-api:8081

---

## 5. Pipelines de desarrollo — 15%
### 5.1 Build de imágenes (matricial)
- Workflow: `.github/workflows/build-images.yml`.
- Disparadores:
  - `push` a `develop` y `master`.
  - `pull_request` a `develop`.
- Estrategia:
  - **Matriz** para `frontend`, `users-api`, `todos-api`, `auth-api`, `log-message-processor`.
  - `docker buildx build` + `push` a Docker Hub: `santi1761/<servicio>`.
  - Tags consistentes (ejemplo): `dev`, `dev-<sha>` en `develop`; `prod`, `prod-<sha>` en `master`; siempre `sha`.
- **Evidencia**: runs verdes y tags visibles en Docker Hub.

### 5.2 Calidad básica / Smoke
- Tras deploy, smoke tests: `curl` a `/`, `/users/johnd`, login y `/todos` con token.

---

## 6. Pipelines de infraestructura — 5%
- `infra-plan.yml`: genera plan de Terraform (se publica en logs).
- `infra-apply-staging.yml` / `infra-apply-prod.yml`: aplican cambios.
- Módulos Terraform:
  - `resource_group/`, `network/`, `keypair/`, `compute_linux_vm/`.
- Variables por ambiente en `infra/terraform/vm/envs/*/terraform.tfvars`.
- **Evidencia**: captura de `Apply complete!` o `No changes`.

---

## 7. Implementación / Despliegue — 20%
### 7.1 Environments y secretos
- **Repository secrets**: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN` (PAT Read & Write con expiración).
- **Environment secrets**:
  - `staging`: `SSH_HOST`, `SSH_USER`, `SSH_PRIVATE_KEY`, `JWT_SECRET`.
  - `production`: `SSH_HOST`, `SSH_USER`, `SSH_PRIVATE_KEY`, `JWT_SECRET`.
- **Protecciones**:
  - `production` con **Required reviewers** (gating).

### 7.2 Deploy STAGING
- Workflow: `.github/workflows/deploy-staging.yml`.
- Conecta por SSH (appleboy/ssh-action), renderiza `.env` desde `ops/compose/env.j2`, sube stack con `docker compose` en `/opt/microservice-app-stg`.
- Nginx expuesto en `:8080`.

### 7.3 Deploy PRODUCCIÓN
- Workflow: `.github/workflows/deploy-prod.yml`.
- Dispara con `push` a `master` o manual (`workflow_dispatch`).
- Gating: requiere aprobación del **environment production**.
- Sube stack en `/opt/microservice-app`, Nginx en `:80`.
- Incluye **smoke tests**: `/`, `/users/johnd`, `/login`, `/todos` con token.

### 7.4 Nginx (reverse proxy)
- Archivo: `ops/compose/nginx.conf` (incluido).
- Reglas: ver ruteo en §4.

---

## 8. Cómo reproducir (paso a paso)
### 8.1 Build de imágenes
```bash
# Commit a develop → dispara build-images (matriz)
# Ver resultados en Actions y en Docker Hub
```

### 8.2 Deploy a STAGING (manual o pull/merge a develop)
```bash
# En Actions: Run workflow (deploy-staging) o push a develop
# Ver /opt/microservice-app-stg en la VM
curl -I http://<IP>:8080/
```

### 8.3 Probar Cache-Aside y Circuit Breaker (STAGING)
```bash
# Reset caché
docker exec -i microservice-app-stg-redis-1 redis-cli FLUSHALL

# Token
TOKEN=$(curl -s -X POST 'http://<IP>:8080/login' \
  -H 'Content-Type: application/json' \
  -d '{"username":"johnd","password":"foo"}' | jq -r '.accessToken // .token')

# Calentar caché
curl -s -H "Authorization: Bearer $TOKEN" http://<IP>:8080/api/todos/ > /dev/null

# Ver claves
docker exec -i microservice-app-stg-redis-1 redis-cli KEYS \*

# Simular falla Users
ssh <vm> 'docker compose -p microservice-app-stg stop users-api'
curl -i -H "Authorization: Bearer $TOKEN" http://<IP>:8080/api/todos/
ssh <vm> 'docker compose -p microservice-app-stg start users-api'
```

### 8.4 Deploy a PRODUCCIÓN (si aplica)
```bash
# Merge develop → master o Run workflow (deploy-prod)
curl -I http://<IP>/
```

---

## 9. Convenciones de tagging de imágenes
- `develop`: `dev`, `dev-<sha>`, `sha`.
- `master`: `prod`, `prod-<sha>`, `sha`.
- Los despliegues leen `TAG` desde el pipeline (`IMAGE_TAG`), por defecto `latest` o el seleccionado en `workflow_dispatch`.

---

## 10. Seguridad y tokens
- **Docker Hub PAT**: permisos **Read & Write**, con expiración (30–90 días) y rotación.  
  1) Generar nuevo PAT → actualizar `DOCKERHUB_TOKEN` → revocar anterior.  
- **SSH**: claves por environment; nunca en el repo.  
- **JWT_SECRET**: secreto por environment.

---

## 11. Troubleshooting
- **“port is already allocated”**: baja el stack del otro ambiente o usa puertos distintos por env y `project_name` diferente en Compose.  
- **401 Unauthorized**: falta header `Authorization: Bearer <token>`.  
- **No aparecen claves en Redis**: asegúrate de hacer una llamada válida (con token) a `/api/todos/` y revisa el Redis del ambiente correcto.  
- **CI no dispara en prod**: confirma que la rama es `master`, que el workflow incluye `branches: [ master ]`, y que el environment `production` está aprobado.

---

## 12. Anexos
- **PlantUML**: ver `docs/arquitectura.puml` (o el archivo `arquitectura.puml` adjunto a este README).  
- **Nginx**: ver `ops/compose/nginx.conf`.  
- **Terraform**: módulos en `infra/terraform/vm/modules/*` y variables en `infra/terraform/vm/envs/*/terraform.tfvars`.

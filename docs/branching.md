# Estrategias de Branching (Dev y Ops) + CI/CD

**Archivo:** `docs/branching.md`  
**Objetivo:** Cumplir los puntos 1 y 2 del enunciado (estrategia de branching para desarrolladores y para operaciones) y dejar el mapa claro de qué pipelines corren en cada rama.

---

## 0) Objetivos

- Separar desarrollo de operaciones para no romper prod.  
- Establecer un flujo claro desde `feature → develop → release → main`.  
- Despliegues por entorno gobernados por ramas `env/` (dev, staging, prod).  
- Etiquetado de imágenes y versiones consistente (SemVer).  

---

## 1) Estrategia para Desarrolladores

### Ramas
- `feature/*` → trabajo diario (una por tarea/issue).  
- `develop` → integración continua de features.  
- `release/x.y.z` → estabilización previa a producción.  
- `main` → producción (**protegida; solo entra código estable**).  

### Reglas de protección (GitHub)
- **main:** requerir PR, 1+ review, checks verdes, impedir push directo.  
- **develop:** opcionalmente proteger igual (recomendado).  

### Flujo

**Crear una feature**
```bash
git checkout -b feature/<resumen-tarea> develop
# ... commits ...
git push -u origin feature/<resumen-tarea>
```
➡️ PR → `develop` (CI build; imágenes tag: `dev-<shortSHA>`).

**Planear una release**
```bash
git checkout -b release/1.0.0 develop
git push -u origin release/1.0.0
```
➡️ Bugs se corrigen allí. CI etiqueta imágenes como `rc-1.0.0`.

**Publicar**
```bash
# merge release -> main
git tag -a 1.0.0 -m "Release 1.0.0"
git push origin 1.0.0
```
➡️ CI de `main` publica imágenes con tag `1.0.0`.

**Back-merge**
```bash
# merge release -> develop
```

**Hotfix en producción**
```bash
git checkout -b hotfix/1.0.1 main
# fix...
git push -u origin hotfix/1.0.1
```
➡️ PR → `main`, crear tag `1.0.1`, luego merge también a `develop`.

---

## 2) Estrategia para Operaciones

La infraestructura y despliegue por entorno se gobiernan con ramas `env/*`.

### Ramas
- `env/dev` → despliega a desarrollo.  
- `env/staging` → despliega a staging (pre-prod).  
- `env/prod` → despliega a producción.  

### Contenido versionado en `env/*`
- `ops/ansible/inventory/<env>.ini` → IP/usuario/puerto SSH de la VM.  
- `ops/compose/.env.j2` → plantilla con `DOCKERHUB_USER` y `TAG`.  
- `ops/compose/docker-compose.yml` → (sin `version:`; Compose v2).  
- `ops/compose/nginx.conf` → rutas correctas (todos 8082, users 8083, auth 8081).  
- *(Opcional infra)* `vm/envs/<env>/terraform.tfvars`.  

>  Archivos de `ops/compose/` son iguales entre entornos; las diferencias llegan por **variables/plantillas en el pipeline**.

### Flujo de cambios desde desarrollo
- `develop` ➜ merge/cherry-pick de `ops/**` hacia `env/dev`.  
- `release/x.y.z` ➜ merge/cherry-pick hacia `env/staging`.  
- `main` (tag `x.y.z`) ➜ merge/cherry-pick hacia `env/prod`.  

---

## 3) Mapa de Pipelines

| Rama         | Qué dispara                           | Qué hace                                                                 |
|--------------|---------------------------------------|--------------------------------------------------------------------------|
| feature/*    | Build (por microservicio)             | Construye y publica imágenes `dev-<shortSHA>`                           |
| develop      | Build (por microservicio)             | Construye y publica imágenes `dev-<shortSHA>`                           |
| release/x.y.z| Build                                 | Construye y publica imágenes `rc-x.y.z`                                 |
| main         | Build + (opcional) Release            | Construye y publica imágenes `x.y.z`; crea tag Git `x.y.z`              |
| env/dev      | `deploy-dev.yml` (+ infra opcional)   | Renderiza `.env`, sube compose/nginx, despliega en VM dev               |
| env/staging  | `deploy-staging.yml` (+ infra opcional)| Igual a dev pero con `rc-x.y.z` y VM staging                            |
| env/prod     | `deploy-prod.yml` (+ infra opcional)  | Igual a staging pero con `x.y.z` y VM producción                        |

**Tagging de imágenes**
- `develop / feature/*` → `dev-<shortSHA>`  
- `release/x.y.z` → `rc-x.y.z`  
- `main (tag x.y.z)` → `x.y.z`  

---

## 4) Secrets y configuración

### Comunes (build)
- `DOCKERHUB_USERNAME`  
- `DOCKERHUB_TOKEN`  

### Deploy por entorno
- **Dev:** `DEV_HOST`, `DEV_USER`, `DEV_SSH_KEY`  
- **Staging:** `STG_HOST`, `STG_USER`, `STG_SSH_KEY`  
- **Prod:** `PROD_HOST`, `PROD_USER`, `PROD_SSH_KEY`  

> `*_SSH_KEY` es clave privada RSA (PEM) con acceso a la VM.  

### Terraform en CI (opcional)
- `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_TENANT_ID`, `ARM_SUBSCRIPTION_ID`  

---

## 5) Guía paso-a-paso

**Desarrollar**
```bash
git checkout -b feature/mejora-x develop
git push -u origin feature/mejora-x
# PR -> develop (se construyen imágenes dev-<sha>)
```

**Subir a dev**
- Merge/cherry-pick de `ops/**` a `env/dev`.  
- Push a `env/dev` ➜ corre `deploy-dev`.  

**Preparar release**
```bash
git checkout -b release/1.0.0 develop
git push -u origin release/1.0.0
```
- CI construye imágenes `rc-1.0.0`.  
- Merge/cherry-pick de `ops/**` a `env/staging`.  

**Publicar**
- Merge `release/1.0.0` → `main`.  
- Crear tag `1.0.0`.  
- Merge/cherry-pick a `env/prod` (usa TAG `1.0.0`).  

**Back-merge**
- Merge `release/1.0.0` → `develop`.  

**Hotfix**
- `hotfix/1.0.1` desde `main` → tag `1.0.1` → `env/prod` → back-merge a `develop`.  

---

## 6) Normas para ops/compose y Nginx

- No usar `version:` en `docker-compose.yml`.  
- Rutas correctas en `nginx.conf`:
  ```
  /api/todos/ → http://todos-api:8082/
  /api/users/ → http://users-api:8083/
  /api/auth/  → http://auth-api:8081/
  ```
- SPA con `try_files`:
  ```nginx
  location / {
    root /usr/share/nginx/html;
    index index.html;
    try_files $uri /index.html;
  }
  ```

---

## 7) Checklist de cumplimiento (Taller)

- [x] Estrategia de branching Dev (`feature/develop/release/main`)  
- [x] Estrategia de branching Ops (`env/dev`, `env/staging`, `env/prod`)  
- [x] Tagging de imágenes (`dev-sha`, `rc-x.y.z`, `x.y.z`)  
- [x] Pipelines de build por microservicio  
- [x] Pipelines de deploy por entorno (Ansible + Docker Compose)  
- [x] *(Opcional)* Pipelines de infraestructura con Terraform  
- [x] NGINX/compose alineados (evita 502)  
- [x] Documentación (este archivo + `docs/arquitectura.md`, `docs/infra.md`, `docs/cicd.md`, `docs/runbook.md`, `docs/known-issues.md`)  

---

## 8) Notas útiles

- **Users API (Spring Boot 1.5)** → usar Java 8.  
- **log-message-processor** → usa `redis==2.10.6`; debe construirse con Python 3.11 (no 3.12).  
- Diferencias por entorno **no se hardcodean** en compose/nginx, se pasan vía `.env` y Ansible.  

---

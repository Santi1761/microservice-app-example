# Branching – Desarrollo (Aplicaciones)

> Objetivo: separar trabajo en curso, integración y releases, y atar cada rama a sus pipelines.

## Ramas
- `develop` (por defecto): integración continua. Dispara builds + **deploy a staging**.
- `feature/*`: ramas cortas desde `develop`. PR -> `develop` con checks obligatorios.
- `release/x.y.z`: estabilización antes de salir. Se despliega a **staging**; si está OK -> merge a `main` + tag `vX.Y.Z` + back-merge a `develop`.
- `main` (o `master`): solo recibe merges desde `release/*` o `hotfix/*`. Dispara **deploy a producción**.
- `hotfix/x.y.z`: desde `main` para urgencias; al cerrar, se back-mergea a `develop`.

## Reglas de PR
- A `develop`: 1 aprobación + status checks verdes (build de servicios modificados).
- A `main`: solo desde `release/*` o `hotfix/*`, con tag de versión y aprobación.

## Versionado de imágenes
- Siempre publicamos `:latest` y `:<git-sha>`.
- En release/tag `vX.Y.Z` también se publica `:vX.Y.Z` (opcional).

## Flujo típico
1. `git checkout -b feature/ajuste-navbar`
2. Commit → PR a `develop` → build + push a Docker Hub → auto-deploy a *staging*.
3. Crear `release/1.0.0` → smoke tests en *staging*.
4. Merge `release/1.0.0` → `main` + tag `v1.0.0` → deploy a *producción*.
5. Back-merge `main` → `develop`.

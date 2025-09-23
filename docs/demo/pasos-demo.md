# Guion de Demo

1. **Cambio mínimo** en `frontend` (texto visible: “Taller 1 – <Nombre> vX”).
2. Commit en `feature/ajuste-texto` → **PR a `develop`**.
3. Ver en **Actions**: builds verdes y **deploy a staging** (environment `staging`).
4. Crear **`release/1.0.0`** desde `develop` → confirma staging OK.
5. Merge `release/1.0.0` → **`main`** + **tag `v1.0.0`** → **deploy a producción** (environment `production` con aprobación).
6. Capturas: pipelines verdes, imágenes en Docker Hub, contenedores corriendo en la VM.

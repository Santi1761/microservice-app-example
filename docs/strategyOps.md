# Branching – Operaciones (Infraestructura)

> Objetivo: controlar cambios de Terraform y despliegue de infraestructura.

## Ramas
- `infra/develop`: trabajo en curso. Corre `terraform fmt/validate/plan` (solo **plan**).
- `infra/staging`: lista para probar. `terraform apply` a *staging* (si hubiera infra separada).
- `infra/main`: aprobada para producción. `terraform apply` a *production* (aprobación manual).

> Nota: Si hoy usamos una sola VM, `infra/staging` puede quedar con “plan only” y aplicamos en `infra/main`.

## Credenciales y secretos
- En **Environments** de GitHub: `staging` y `production`.
- Secrets mínimos por env: 
  - `SSH_HOST`, `SSH_USER`, `SSH_PRIVATE_KEY` (para desplegar con Ansible/SSH).
  - Si Terraform en Azure vía SP: `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_TENANT_ID`, `ARM_SUBSCRIPTION_ID`.

## Flujo típico
1. PR a `infra/develop` → `infra-plan.yml` genera el **plan** (sin aplicar).
2. Merge a `infra/main` (aprobación) → `infra-apply-prod.yml` ejecuta `terraform apply`.

# Pipelines – Infraestructura (Terraform)

Workflows:
- `infra-plan.yml`: corre en PR/push a `infra/develop` (y similares), ejecuta `init`, `fmt -check`, `validate`, `plan`.
- `infra-apply-staging.yml`: corre en `push` a `infra/staging` (si hay entorno staging), ejecuta `apply`.
- `infra-apply-prod.yml`: corre en `push` a `infra/main` con **environment: production** y aprobación, ejecuta `apply`.

Secretos (si aplica Azure SP):
- `ARM_CLIENT_ID`, `ARM_CLIENT_SECRET`, `ARM_TENANT_ID`, `ARM_SUBSCRIPTION_ID`.

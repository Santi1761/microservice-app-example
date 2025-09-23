David Santiago Donneys
Santiago Arboleda Velasco

# Taller 1 – Automatización de Infraestructura con Pipelines

Repositorio: `microservice-app-example` (fork).  
Objetivo: Construir **pipelines de infraestructura y de código** que permitan a un equipo ágil desarrollar, integrar y desplegar los 5 servicios del proyecto.

## Qué contiene este repositorio
- **Aplicaciones**: `auth-api` (Go), `users-api` (Spring Boot), `todos-api` (Node), `log-message-processor` (Python), `frontend` (Vue).
- **Infraestructura**: `vm/` (Terraform) para Azure VM y red; `ops/compose/` con `docker-compose` para staging/prod; `ops/ansible/` para despliegue automatizado en la VM.
- **CI/CD**: Workflows en `.github/workflows/` para build por servicio, deploy a *staging* y *producción*, y plan/apply de Terraform.
- **Documentación** Metodología, branching, pipelines, arquitectura y guion de demo.

## Navegación rápida
- Metodología: [/docs/methodology.md](./methodology.md)
- Branching Dev: [/docs/branching-dev.md](./docs/strategyDev.md)
- Branching Ops/Infra: [/docs/branching-ops.md](./docs/strategyOps.md)
- Pipelines (build/test/push): [/docs/pipelines/dev-workflows.md](./pipelines/dev-workflows.md)
- Pipelines (deploy): [/docs/pipelines/deploy-workflows.md](./pipelines/deploy-workflows.md)
- Pipelines (infra): [/docs/pipelines/infra-workflows.md](./pipelines/infra-workflows.md)
- Arquitectura (diagrama y decisiones): [/docs/architecture/arquitectura.md](./architecture/arquitectura.md)
- Guion de demo: [/docs/demo/pasos-demo.md](./demo/pasos-demo.md)

## Entregables principales (según rúbrica)
1. Estrategia de branching (dev y ops).
2. ≥ 2 patrones de diseño de nube (documentados en arquitectura).
3. Diagrama de arquitectura actualizado al repo.
4. Pipelines de desarrollo y de infraestructura.
5. Implementación en Azure VM y evidencia de despliegue.
6. Presentación/demostración: cambio en `develop` → staging; release → `main` → producción. 
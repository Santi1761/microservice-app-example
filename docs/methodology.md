# Metodología Ágil

Elegimos un **Scrum-ligero** (1 semana por iteración) centrado en habilitar CI/CD:

- **Roles**: 
  - Dev Team (equipo): desarrollo microservicios y pruebas.
  - Ops (equipo): infraestructura y despliegue.
  - Stakeholder (profesor).

- **Ceremonias**:
  - **Planning**: objetivo del sprint = “pipeline vivo de build y despliegue por servicio”.
  - **Daily**: coordinación corta (chat del equipo).
  - **Review**: demo del pipeline corriendo y despliegue a staging/prod.
  - **Retro**: mejoras (p. ej. cache de dependencias, tiempos de build).
  
- **Definición de Hecho (DoD)**:
  - Código en rama correcta (según branching).
  - Build verde en Actions y **imagen publicada** en Docker Hub.
  - Despliegue a ambiente correspondiente (staging o prod).
  - Evidencias añadidas en `/docs/demo/evidencias/`.

- **Gestión de trabajo**:
  - Issues/PRs en GitHub, convención de ramas (`feature/*`, `release/*`, `hotfix/*`).

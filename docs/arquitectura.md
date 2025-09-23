# Arquitectura de la Solución

## Visión general
Cliente (web) → Nginx (reverse proxy) → Microservicios:
- **frontend** (Vue)
- **auth-api** (Go)
- **users-api** (Spring Boot)
- **todos-api** (Node)
- **log-message-processor** (Python)

Infraestructura:
- **Azure VM** (Docker + Docker Compose).
- **Redis** (para soporte de caché y colas de logs).
- Red interna de Docker (`micro-net`).

> Ver diagrama: `architecture/diagram.png` (editable `diagram.drawio`). Puedes partir de `arch-img/Microservices.png` y añadir: VM Azure, Nginx, Redis, puertos y dependencias.

## Decisiones relevantes
- **Branching**: `feature/*` → `develop` (staging) → `release/*` → `main` (producción).  
- **CI/CD**: build por microservicio; despliegue por environment con aprobación en producción.
- **Patrones de nube documentados**:
  - *Cache-Aside* (capa de cache Redis para endpoints de lectura frecuente).
  - *Circuit Breaker* (degradación controlada si `users-api` no responde).
- **Observabilidad**: `log-message-processor` para centralizar logs.

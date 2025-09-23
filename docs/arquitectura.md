# Arquitectura del sistema (microservice-app-example)

## Vista lógica
- **Frontend (Nginx)** expuesto en **:80**: sirve la SPA y actúa como *reverse proxy* hacia:
  - `/api/todos` → **todos-api (Node.js)** – Puerto interno **8082**
  - `/api/users` → **users-api (Spring Boot 1.5)** – Puerto interno **8083**
  - `/api/auth`  → **auth-api (Go)** – Puerto **8081**
- **Redis**: cola usada por todos-api para publicar eventos (CREATE/DELETE).
- **log-message-processor (Python)**: consumidor de Redis que procesa los mensajes.

> Imagen de referencia: `arch-img/Microservices.png`

## Vista de despliegue
- **Azure**:
  - **Resource Group**, **Virtual Network**, **Subnet**, **Network Security Group**, **Public IP**, **VM Ubuntu 22.04**.
- **En la VM**:
  - Docker + Compose v2 instalados por **Ansible**.
  - Carpeta de despliegue: `/opt/microservice-app`
  - Arranque con `docker compose up -d` usando `docker-compose.yml`, `.env` y `nginx.conf`.

## Patrones de nube implementados
1. **Cache-Aside**: 
   - Primera lectura de `/api/todos` va a la fuente y se guarda en cache; la segunda lectura responde con `fromCache: true`.
   - Demostración: invocar dos veces `GET /api/todos` consecutivas; la segunda respuesta indica `fromCache`.
2. **Circuit Breaker**:
   - Si **users-api** no responde, **todos-api** utiliza un **fallback** (p. ej. “Unknown User”) sin tumbar la solicitud.
   - Demostración: `docker compose stop users-api`; consumir `/api/todos`; se observa la respuesta con fallback; luego `docker compose start users-api`.

## Puertos y ruteo (clave para evitar 502)
- **Puertos REALES (dentro de contenedor)**: `todos-api:8082`, `users-api:8083`, `auth-api:8081`.
- **Nginx** apunta exactamente a esos puertos internos (no 3000/8080).
- Soporte con y sin “/” final: `/api/todos` redirige a `/api/todos/`.

## Seguridad básica
- NSG con **22** (SSH) y **80** (HTTP).
- Acceso SSH con clave **RSA 4096** (Azure no acepta ed25519).

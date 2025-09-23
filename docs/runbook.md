# Runbook (operación)

## VM
ssh -i private_ssh_key.pem ubuntu@$(terraform output -raw public_ip)
cd /opt/microservice-app

## Estado y logs
docker compose ps
docker compose logs -n 80 frontend todos-api users-api log-processor

## Ver config de Nginx cargada
docker compose exec frontend sh -lc 'cat /etc/nginx/conf.d/default.conf'

## Reinicios
docker compose restart <servicio>
docker compose up -d

## Demo Cache-Aside
# Desde tu máquina
Invoke-RestMethod -Uri "http://PUBLIC_IP/api/todos"
Invoke-RestMethod -Uri "http://PUBLIC_IP/api/todos" | Format-List

## Demo Circuit Breaker
# En la VM
docker compose stop users-api
# Desde tu máquina llamar a la operación que requiere user -> ver fallback
docker compose start users-api

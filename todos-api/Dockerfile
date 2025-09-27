# todos-api/Dockerfile
FROM node:20-alpine
WORKDIR /app

COPY package*.json ./

RUN npm install --omit=dev --no-audit --no-fund

# Resto del código
COPY . .

ENV NODE_ENV=production PORT=8082
EXPOSE 8082

CMD ["node", "server.js"]

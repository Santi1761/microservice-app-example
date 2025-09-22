
FROM node:8.17.0 AS build
WORKDIR /app

# Dependencias de compilación para node-sass si no hay binario precompilado
RUN apt-get update && apt-get install -y python make g++ && rm -rf /var/lib/apt/lists/*

COPY package.json ./
COPY package-lock.json ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80

FROM mirror.gcr.io/library/node:20-alpine
WORKDIR /app

COPY package*.json ./
RUN (npm ci --omit=dev --no-audit --no-fund) || \
    (npm install --omit=dev --no-audit --no-fund)

COPY . .
ENV NODE_ENV=production PORT=8082
EXPOSE 8082
CMD ["node", "server.js"]

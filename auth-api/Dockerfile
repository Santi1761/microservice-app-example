FROM golang:1.22 AS build
WORKDIR /src
COPY . .

# Si no hay go.mod, lo crea; luego resuelve dependencias
RUN if [ ! -f go.mod ]; then go mod init auth-api; fi \
    && go mod tidy

# Compilar binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/auth ./...

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/auth /auth
EXPOSE 8081
ENTRYPOINT ["/auth"]

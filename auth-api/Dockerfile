# Builder (golang) desde el mirror de Google (no Docker Hub)
FROM mirror.gcr.io/library/golang:1.22 AS build
WORKDIR /src
COPY . .

# Si no hay go.mod, lo crea; luego deps
RUN if [ ! -f go.mod ]; then go mod init auth-api; fi \
    && go mod tidy

# Binario estático
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/auth ./...

# Runtime distroless (GCR oficial de Google)
FROM gcr.io/distroless/base-debian12
COPY --from=build /out/auth /auth
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/auth"]

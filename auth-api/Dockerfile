FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/auth

FROM gcr.io/distroless/base-debian12
COPY --from=build /out/auth /auth
EXPOSE 8081
ENTRYPOINT ["/auth"]
# Build stage
FROM golang:1.23.0-bullseye AS build-env

# Installer UPX
RUN apt-get update && apt-get install -y upx

WORKDIR /app

COPY . .

RUN go mod download

# Compilation et optimisation du code Go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/main

# Compression de l'exécutable avec UPX
RUN upx --best /app/main

# Nettoyage des caches et fichiers temporaires
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Deployment stage
FROM gcr.io/distroless/static-debian11

COPY --from=build-env /app/main /app/main

CMD ["/app/main"]
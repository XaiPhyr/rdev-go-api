# STAGE 1: Build
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache build-base
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate

# STAGE 2: Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binaries and config
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/internal/data ./internal/data
COPY --from=builder /app/config.docker.yaml ./config.yaml

# Copy the entrypoint script from your local scripts folder
COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8200

ENTRYPOINT ["/entrypoint.sh"]
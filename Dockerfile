# STAGE 1: Build
FROM golang:1.26-alpine AS builder

# Install build-base for any potential C dependencies (like some SQLite drivers)
RUN apk add --no-cache build-base

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/api

# STAGE 2: Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/main .
COPY --from=builder /app/config.docker.yaml ./config.yaml 

EXPOSE 8200

CMD ["./main"]
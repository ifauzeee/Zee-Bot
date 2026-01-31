# Build Stage
FROM golang:1.25-alpine AS builder

# Install build dependencies first (cached layer)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
# Cache go modules
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# Cache build artifacts
# CGO_ENABLED=1 required for go-sqlite3
ENV CGO_ENABLED=1

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-w -s" -o zee-ubot ./cmd/bot

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install basic dependencies (ca-certificates for HTTPS)
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/zee-ubot .

# Volume for session
VOLUME /app/session

CMD ["./zee-ubot"]

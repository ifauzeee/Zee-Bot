# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
# Cache go modules
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# CGO_ENABLED=0 makes build super fast (Turbo) and the binary static
ENV CGO_ENABLED=0

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

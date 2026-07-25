# ── Stage 1: Build ─────────────────────────────────────────────────────────
# Use the full Go image so we have GCC available for CGO (sqlite3).
FROM golang:1.23-bookworm AS builder

WORKDIR /app

# Download dependencies first (cached layer — only re-runs if go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Install templ code generator
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy full source
COPY . .

# Generate Templ templates
RUN templ generate

# Build static binary with version injected
ARG VERSION=dev
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-s -w -X main.version=${VERSION} -X github.com/Suke2004/atlas-go/internal/health.version=${VERSION}" \
    -o /app/build/atlas \
    ./cmd/atlas

# ── Stage 2: Runtime ───────────────────────────────────────────────────────
# Debian slim gives us a minimal runtime with libc (needed by CGO/sqlite3).
FROM debian:bookworm-slim

# Install CA certificates (for HTTPS calls to AI providers) and sqlite3 tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user for the process
RUN useradd -r -s /bin/false -m -d /app atlas

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/build/atlas ./atlas

# The /data directory is always a mounted volume — never baked into the image.
# We create the directory here so Docker volume mount works correctly.
RUN mkdir -p /app/data/db /app/data/uploads /app/data/backups /app/data/logs \
    && chown -R atlas:atlas /app

USER atlas

EXPOSE 8080

# Health check — Docker will restart the container if this fails
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./atlas"]
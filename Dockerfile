# syntax=docker/dockerfile:1

# ---- Stage 1: build frontend (landing + admin, single Vite project) ----
FROM oven/bun:1-alpine AS frontend-builder
# vue-tsc patches Node's fs.readFileSync to resolve .vue files; Bun's runtime
# doesn't support that hook, so it needs a real node to run under.
# ponytail: revisit once oven-sh/bun#4754 ships a bun-native vue-tsc.
RUN apk add --no-cache nodejs
WORKDIR /build
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend/ ./
RUN bun run build

# ---- Stage 2: build Go backend ----
FROM golang:1.24-alpine AS backend-builder
WORKDIR /build
RUN apk add --no-cache git
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Embed unified frontend dist (landing at /, admin at /admin/) so Go serves both
RUN mkdir -p public
COPY --from=frontend-builder /build/dist/ ./public/
# ncruces/go-sqlite3 is pure Go (WASM) -> no CGO needed
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/okuru .

# ---- Stage 3: minimal runtime ----
FROM alpine:3.20 AS runtime
RUN apk add --no-cache ca-certificates tzdata wget postgresql16-client \
    && addgroup -S okuru && adduser -S okuru -G okuru
WORKDIR /opt/okuru

# Application binary
COPY --from=backend-builder /out/okuru ./okuru
# Runtime assets required by Goravel on disk
COPY backend/config/        ./config/
COPY backend/resources/     ./resources/
COPY backend/database/      ./database/
COPY --from=backend-builder /build/public/ ./public/
COPY scripts/               ./scripts/

# Artisan wrapper: avoids needing the Go toolchain at runtime.
RUN printf '#!/bin/sh\nexec /opt/okuru/okuru artisan "$@"\n' > ./artisan && chmod +x ./artisan \
    && mkdir -p storage/app/public storage/framework/cache/data storage/logs \
    && chown -R okuru:okuru /opt/okuru

USER okuru
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q --spider http://127.0.0.1:8080/health || exit 1

CMD ["./okuru"]

# syntax=docker/dockerfile:1

# ── Stage 1: Build Tailwind CSS ───────────────────────────────────────────────
FROM node:22-alpine AS css-builder
WORKDIR /build

# Install deps first (cached unless package files change)
COPY package.json package-lock.json ./
RUN npm ci

# Copy only the files Tailwind needs to scan for class names
COPY static/css/input.css static/css/input.css
COPY guides/     guides/
COPY templates/  templates/
COPY handlers/   handlers/
COPY main.go     ./

RUN ./node_modules/.bin/tailwindcss \
    -i static/css/input.css \
    -o static/css/output.css \
    --minify


# ── Stage 2: Generate templ + compile Go binary ───────────────────────────────
FROM golang:1.26-alpine AS go-builder
WORKDIR /build

# Install templ generator — pinned to match go.mod
RUN go install github.com/a-h/templ/cmd/templ@v0.3.1001

# Download modules (cached unless go.mod/go.sum change)
COPY go.mod go.sum ./
RUN go mod download

# Copy full source and build
COPY . .
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /out/stylesheets \
    .


# ── Stage 3: Minimal runtime image ────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12 AS runtime
WORKDIR /app

# Binary
COPY --from=go-builder /out/stylesheets /app/stylesheets

# Static assets — copy directory first, then overwrite CSS with built output
COPY --from=go-builder /build/static /app/static
COPY --from=css-builder /build/static/css/output.css /app/static/css/output.css

EXPOSE 8080
USER nonroot:nonroot
CMD ["/app/stylesheets"]

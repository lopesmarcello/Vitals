# Stage 1: Build
# We explicitly pin Go 1.23 to satisfy the templ requirement
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate templates
RUN templ generate

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o vitals ./cmd/server/main.go

# Stage 2: Runtime
FROM alpine:latest

# Install Chromium and dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    dumb-init

# Environment variables for Chrome
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

WORKDIR /app
COPY --from=builder /app/vitals .

# Use dumb-init
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# Run the app
CMD ["./vitals"]

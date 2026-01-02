# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder
WORKDIR /app

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templates
RUN templ generate

# Build the binary
RUN go build -o vitals ./cmd/server/main.go

# Stage 2: The Runtime Image (Alpine + Chromium)
FROM alpine:latest

# Install Chromium and dependencies
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Tell Chromedp where to find Chrome
ENV CHROME_BIN=/usr/bin/chromium-browser
ENV CHROME_PATH=/usr/lib/chromium/

WORKDIR /app
COPY --from=builder /app/vitals .

# Expose the port
EXPOSE 8080

# Run it
CMD ["./vitals"]

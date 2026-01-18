# Use an official Go image as the build stage
FROM golang:tip-bookworm AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary with version information
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=${VERSION}" -o subscription-service

# Use a minimal base image for final container
FROM debian:bullseye-slim

# Install Chrome/Chromium dependencies and CA certificates
RUN apt-get update && apt-get install -y \
    ca-certificates \
    dumb-init \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/subscription-service .

# Create /dev/shm volume directory
RUN mkdir -p /dev/shm

# Set environment variables (can be overridden at runtime)
ENV REDIS_HOST=redis-service \
    REDIS_PORT=6379 \
    REDIS_PASSWORD= \
    RABBITMQ_HOST=rabbitmq-service \
    RABBITMQ_PORT=5672 \
    RABBITMQ_USERNAME=guest \
    RABBITMQ_PASSWORD=guest \
    SMTP_HOST=smtp.local.dev \
    SMTP_PORT=1025
# Expose health check port
EXPOSE 8080

# Use dumb-init as entrypoint to handle signals properly
ENTRYPOINT ["/usr/bin/dumb-init", "--"]

# Command to run the service
CMD ["./subscription-service"]

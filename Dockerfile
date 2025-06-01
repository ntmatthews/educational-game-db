# Educational Game Database Dockerfile
FROM golang:1.21-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o educational-game-db cmd/cli/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and sqlite
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/educational-game-db .

# Copy web assets
COPY --from=builder /app/web ./web

# Create data directory
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV DB_PATH=/data/accounts.db

# Run the application
CMD ["./educational-game-db", "web", "--port", "8080", "--db", "/data/accounts.db"]

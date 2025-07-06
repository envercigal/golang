# 1. Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd

# 2. Final stage
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the built binary
COPY --from=builder /app/server .

# Expose application port
EXPOSE 3000

# Entrypoint
ENTRYPOINT ["./server"]

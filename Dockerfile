# Build stage
FROM docker.io/library/golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o receiver main.go

# Final stage
FROM docker.io/library/alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/receiver .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./receiver"]

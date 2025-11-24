# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS calls if needed
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port 8081
EXPOSE 8081

# Run the binary
CMD ["./main"]

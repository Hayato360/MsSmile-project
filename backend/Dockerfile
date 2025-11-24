# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required system dependencies including sqlite-dev
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for SQLite
# Use external sqlite library to avoid musl compatibility issues
RUN CGO_ENABLED=1 GOOS=linux go build -tags "libsqlite3" -o main .

# Final Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates and sqlite runtime library
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port 8081
EXPOSE 8081

# Run the binary
CMD ["./main"]

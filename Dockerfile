# Stage 1: Build
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# -ldflags="-w -s" reduces the size of the binary by removing debug information.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /main ./cmd/server

# Stage 2: Deploy
FROM alpine:latest

# Add ca-certificates to make SSL calls
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /main .
COPY config.yaml .
COPY docs ./docs

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

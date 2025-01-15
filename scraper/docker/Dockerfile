# Stage 1: Build the Go binary
FROM golang:1.20-alpine AS builder

# Set environment variables
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY ../. .

# Build the binary
RUN go build -o scraper main_scheduler.go

# Stage 2: Create a minimal runtime container
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/scraper .

# Environment variable (can be overridden via docker-compose)
ENV DB_CONN_STRING=postgres://postgres:secret@db:5432/clothes_scraper?sslmode=disable

# Command to run the scraper
CMD ["./scraper"]

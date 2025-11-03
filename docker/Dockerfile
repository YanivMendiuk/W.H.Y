# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY ../go.mod ../go.sum ./
RUN go mod download

# Copy the rest of the source
COPY ../cmd ./cmd
COPY ../internal ./internal

# Build the binary
RUN go build -o webhook ./cmd/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/webhook .

EXPOSE 8181

# Run the binary
CMD ["./webhook"]


# ==========================================
# BUILD STAGE
# ==========================================
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked binary
# TODO: Update path to entrypoint if different (e.g. ./cmd/api)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api

# ==========================================
# RUNNER STAGE (Alpine)
# ==========================================
FROM alpine:3.21

# Create non-root group and user
RUN addgroup -S app && adduser -S -G app app

WORKDIR /app

# Copy CA certificates for external HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/app ./app

# Create config directory and set ownership
# TODO: Modify if config directory path changes
RUN mkdir -p /app/pkg/config && \
    chown -R app:app /app

USER app

# TODO: Update exposed port to match your application config
EXPOSE 38600

CMD ["./main"]

# ==========================================
# OPTIONAL: SCRATCH RUNNER STAGE
# ==========================================
# To use scratch (absolute minimum size, no shell):
#
# FROM scratch
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /app/main /main
# # NOTE: Volume mounts in scratch require binary to run as root,
# # or numeric USER if files on host have matching permissions.
# EXPOSE 38600
# ENTRYPOINT ["/main"]

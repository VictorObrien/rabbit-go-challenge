FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod files
COPY go.mod go.sum ./
COPY services/publisher/go.mod services/publisher/go.sum ./services/publisher/

# Download dependencies
RUN cd services/publisher && go mod download

# Copy source code
COPY pkg/ ./pkg/
COPY services/publisher/ ./services/publisher/

# Build
RUN cd services/publisher && CGO_ENABLED=0 GOOS=linux go build -o /publisher main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /publisher .
CMD ["./publisher"]


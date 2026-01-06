FROM golang:1.23-alpine AS builder

# Install gcc for SQLite
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go.mod files
COPY go.mod go.sum ./
COPY services/consumer/go.mod services/consumer/go.sum ./services/consumer/

# Download dependencies
RUN cd services/consumer && go mod download

# Copy source code
COPY pkg/ ./pkg/
COPY services/consumer/ ./services/consumer/

# Build with CGO enabled for SQLite
RUN cd services/consumer && CGO_ENABLED=1 GOOS=linux go build -o /consumer main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /root/
COPY --from=builder /consumer .
CMD ["./consumer"]


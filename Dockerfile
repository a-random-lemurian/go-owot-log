# Stage 1: build
FROM golang:1.24-bullseye AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download -x
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -trimpath -ldflags="-s -w" -o /app/owot-log ./cmd/owot-log

# Stage 2: scratch runtime (includes CA certs from builder)
FROM scratch
COPY --from=builder /app/owot-log /app/owot-log
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /app
ENTRYPOINT ["/app/owot-log","start"]

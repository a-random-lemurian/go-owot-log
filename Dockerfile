# Stage 1: build
FROM golang:1.24-bullseye AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x


COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -v -trimpath -ldflags="-s -w" -o /app/owot-log ./cmd/owot-log

# Stage 2: scratch runtime (includes CA certs from builder)
FROM gcr.io/distroless/base-debian12
EXPOSE 27321
COPY --from=builder /app/owot-log /app/owot-log
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
WORKDIR /app
ENTRYPOINT ["/app/owot-log","start"]

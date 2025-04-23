FROM golang:1.23.4-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY . .

ARG VERSION=dev
LABEL org.opencontainers.image.version=${VERSION}
RUN GOOS=linux go build -ldflags="-s -w -X main.version=${VERSION}" -o /app/zenfeed ./main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && \
    mkdir -p /app/data

COPY --from=builder /app/zenfeed /app/

WORKDIR /app
ENTRYPOINT ["/app/zenfeed"]
CMD ["--config", "/app/config/config.yaml"]
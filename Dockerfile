FROM golang:1.20-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build --ldflags '-linkmode external -extldflags=-static' main.go

FROM scratch AS runtime
LABEL org.opencontainers.image.source "https://github.com/atorrescogollo/aws-local-serverless-proxy"
LABEL org.opencontainers.image.version "v1.0.1"
LABEL org.opencontainers.image.authors "Álvaro Torres Cogollo <atorrescogollo@gmail.com>"
LABEL org.opencontainers.image.licenses GPL-3.0-or-later
LABEL org.opencontainers.image.title "AWS Local Serverless Proxy"
LABEL org.opencontainers.image.description "HTTP API gateway for locally testing AWS lambdas"
LABEL org.opencontainers.image.url "https://github.com/atorrescogollo/aws-local-serverless-proxy"

COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]

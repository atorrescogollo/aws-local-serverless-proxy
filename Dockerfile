FROM golang:1.19-alpine AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN go build --ldflags '-linkmode external -extldflags=-static' main.go

FROM scratch AS runtime
LABEL org.opencontainers.image.source "https://github.com/atorrescogollo/aws-lambda-serverless-proxy"
LABEL org.opencontainers.image.description HTTP API gateway for local AWS lambdas

COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]

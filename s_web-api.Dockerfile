#- Build Stage
FROM golang:1.18.0-alpine3.15 AS builder

WORKDIR /build

# Copy go.mod and go.sum and download the needed modules
COPY ./pkg/web-api/go.mod ./pkg/web-api/go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./api.dev.yaml ./

# Build the services for the given architecture and os
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build \
      -buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw -a -v -o \
      web-api ./cmd/web-api/main.go

#- Run Stage
FROM alpine:3.15

WORKDIR /app

# Copy the web-api service executable, start script ,and api.dev.yaml config file
COPY --from=builder /build/web-api ./
COPY --from=builder /build/api.dev.yaml ./

EXPOSE 8081
CMD ./web-api

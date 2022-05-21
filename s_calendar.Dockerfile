#- Build Stage
FROM golang:1.18.0-alpine3.15 AS builder

WORKDIR /build

# Copy go.mod and go.sum and download the needed modules
COPY go.mod go.sum ./
RUN go mod download

# Install grpc-health-probe
RUN apk add git &&\
    git clone https://github.com/grpc-ecosystem/grpc-health-probe.git

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./api.dev.yaml  ./

# Build the services for the given architecture and os
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build \
    -buildmode=pie -trimpath -modcacherw -a -v -o \
    calendar ./cmd/calendar/main.go
# Build the grpc-health-probe tool
RUN cd ./grpc-health-probe && \
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build \
    -buildmode=pie -trimpath -modcacherw -a -v -o grpc-health-probe .

#- Run Stage
FROM alpine:3.15

WORKDIR /app

# Copy the calendar service executable and the calendar database migrations
COPY --from=builder /build/calendar ./
COPY --from=builder /build/api.dev.yaml ./
COPY --from=builder /build/grpc-health-probe/grpc-health-probe ./

EXPOSE 8082
ENTRYPOINT ./calendar

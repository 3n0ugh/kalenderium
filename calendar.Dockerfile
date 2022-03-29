#- Build Stage
FROM golang:1.18.0-alpine3.15 AS builder

WORKDIR /build

# Copy go.mod and go.sum and download the needed modules
COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./api.dev.yaml  ./

# Build the services for the given architecture and os
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o calendar ./cmd/calendar/main.go

# Install the go-migrate
RUN apk add curl &&\
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

#- Run Stage
FROM alpine:3.15

WORKDIR /app

# Copy the calendar service executable and the calendar database migrations
COPY --from=builder /build/calendar ./
COPY --from=builder /build/pkg/calendar/database/migrations ./calendar-migrations
COPY --from=builder /build/api.dev.yaml ./
COPY --from=builder /build/migrate ./

EXPOSE 8082
CMD ./migrate -path=./calendar-migrations -database="postgres://kalenderium:kartaca2022@postgres/calendar?sslmode=disable" up && \
    ./calendar

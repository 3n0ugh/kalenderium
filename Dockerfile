#- Build Stage
FROM golang:1.18.0-alpine3.15 AS builder

WORKDIR /app

COPY . .

RUN go build -o calendar ./cmd/calendar/main.go
RUN apk --no-cache add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

#- Run Stage
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /app/calendar .
COPY --from=builder /app/migrate ./migrate
COPY api.dev.yaml .
COPY pkg/calendar/database/migrations ./migration

EXPOSE 8080
CMD [ "/app/calendar" ]

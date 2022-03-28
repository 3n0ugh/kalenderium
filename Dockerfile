#- Build Stage
FROM golang:1.18.0-alpine3.15 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg
COPY ./api.dev.yaml .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o calendar ./cmd/calendar/main.go &&\
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o  account ./cmd/account/main.go &&\
    GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o  web-api ./cmd/web-api/main.go
#RUN apk add curl &&\
#    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

#- Run Stage
FROM alpine:3.15

WORKDIR /app

COPY --from=builder /build/calendar .
COPY --from=builder /build/account .
COPY --from=builder /build/web-api .
#COPY --from=builder /app/start.sh .
COPY --from=builder /build/api.dev.yaml .

#RUN chmod +x ./start.sh

EXPOSE 8080
CMD [ "./calendar" ]

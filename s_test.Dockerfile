# Start from golang base image
FROM golang:1.18.0-alpine3.15

# Set the current working directory inside the container
WORKDIR /test

# Copy go.mod, go.sum files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy sources to the working directory
COPY . .

# Run the test suite
CMD CGO_ENABLED=0 go test -v -cover ./pkg/account \
    | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' \
    | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/'' \
    && CGO_ENABLED=0 go test -v -cover ./pkg/calendar \
    | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' \
    | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/'' \
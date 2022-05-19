# ==================================================================================== #
# 										HELPERS                                        #
# ==================================================================================== #
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# 										DEVELOPMENT 								   #
# ==================================================================================== #
# Create new migration for account database
db/migrate/new/account:
	migrate create -seq -ext=.sql -dir=./pkg/account/database/migrations/ create_user_table

# Create new migration for calendar database
db/migrate/new/calendar:
	migrate create -seq -ext=.sql -dir=./pkg/calendar/database/migrations/ create_event_table

# Up the migrations for account database
db/migrate/up/account: confirm
	migrate -path=./pkg/account/database/migrations -database=${ACCOUNT_DB_DSN} up

# Up the migrations for calendar database
db/migrate/up/calendar: confirm
	migrate -path=./pkg/calendar/database/migrations -database=${CALENDAR_DB_DSN} up

# Generate output from calendar proto file
proto/create/calendar: confirm
	cd pkg/calendar/pb; ./proto.sh

# Generate output from account proto file
proto/create/account: confirm
	cd pkg/account/pb; ./proto.sh

# Test Account gRPC service
test/account:
	go test -v -cover ./pkg/account

# Test Calendar gRPC service
test/calendar:
	go test -v -cover ./pkg/calendar

test/all:
	CGO_ENABLED=0 go test -v -cover ./pkg/account && CGO_ENABLED=0 go test -v -cover ./pkg/calendar
# ==================================================================================== #
# 							 	    PRODUCTION								           #
# ==================================================================================== #

# Build the containers according to docker-compose.yaml
docker/build:
	docker compose build

# Run the containers according to docker-compose.yaml
docker/run:
	docker compose up -d

# Build the test container according to s_test.Dockerfile
docker/test/build:
	docker build -t kalenderium/test -f s_test.Dockerfile .

# Run the test container according to s_test.Dockerfile
docker/test/run:
	docker run -it --rm --name kalenderium-test kalenderium/test

# Stop the containers
docker/stop:
	docker compose down

# Install the node_modules
vue/install:
	cd web && yarn install

# Run the vue frontend app
vue/run:
	cd web && yarn serve

# Docker: ready up the config
config/docker:
	sed -i.bak '/use with docker/s/^#//g;/use without docker/s/^/#/g' api.dev.yaml

# Local: ready up the config
config/local:
	sed -i.bak '/use with docker/s/^/#/g;/use without docker/s/^#//g' api.dev.yaml

# Local: start calendar service
local/run/calendar:
	go run ./cmd/calendar

# Local: start account service
local/run/account:
	go run ./cmd/account

# Local: start web-api service
local/run/web-api:
	go run ./cmd/web-api

# ==================================================================================== #
# 									QUALITY CONTROL								       #
# ==================================================================================== #
## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...
## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy -compat=1.17
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor
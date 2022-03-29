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

# ==================================================================================== #
# 							 	    PRODUCTION								           #
# ==================================================================================== #

# Build the containers according to docker-compose.yaml
docker/build:
	docker-compose build

# Run the containers according to docker-compose.yaml
docker/run:
	docker-compose run -d

# Stop the containers
docker/stop:
	docker-compose down

# Run the vue frontend app
vue/run:
	cd web; yarn serve

# Local: start calendar service
local/run/calendar:
	go run ./cmd/calendar

# Local: start account service
local/run/account:
	go run ./cmd/account

# Local: start web-api service
local/run/web-api:
	go run ./cmd/web-api

db/migrate/new:
	#migrate create -seq -ext=.sql -dir=./pkg/account/database/migrations/ create_user_table
	migrate create -seq -ext=.sql -dir=./pkg/calendar/database/migrations/ create_event_table
db/migrate/up:
	#migrate -path=./pkg/account/database/migrations -database=${ACCOUNT_DB_DSN} up
	migrate -path=./pkg/calendar/database/migrations -database=${CALENDAR_DB_DSN} up

protoc/create:
	protoc -I.  "${PB_PATH}/${PROTO_FILE}" --go_out=plugins=grpc:"${PB_PATH}"

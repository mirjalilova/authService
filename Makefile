CURRENT_DIR=$(shell pwd)
APP=template

DBURL="postgres://postgres:feruza1727@localhost:5432/staffer?sslmode=disable"

run:
	go run main.go
init:
	go mod init
	go mod tidy 
	go mod vendor

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

migrate_up:
	migrate -path migrations -database postgres://postgres:feruza1727@localhost:5432/timecapsule?sslmode=disable -verbose up

migrate_down:
	migrate -path migrations -database postgres://postgres:feruza1727@localhost:5432/timecapsule?sslmode=disable -verbose down

migrate_force:
	migrate -path migrations -database postgres://postgres:feruza1727@localhost:5432/timecapsule?sslmode=disable -verbose force 1

migrate_file:
	migrate create -ext sql -dir migrations -seq create_table

build:
	CGO_ENABLED=0 GOOS=darwin go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

swag-gen:
	~/go/bin/swag init -g ./api/api.go -o api/docs force 1
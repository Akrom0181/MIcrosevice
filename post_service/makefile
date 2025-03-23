CURRENT_DIR=$(shell pwd)

APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd

gen-proto-module:
	sudo rm -rf ${CURRENT_DIR}/genproto
	./scripts/gen_proto.sh ${CURRENT_DIR}
	
swag_init:
	swag init -g api/main.go -o api/docs

migrate-up:
	migrate -database 'postgres://akromjonotaboyev:1@localhost:5432/microservice?sslmode=disable' -path migrations up;

migrate-down:
	migrate -database 'postgres://akromjonotaboyev:1@localhost:5432/microservice?sslmode=disable' -path migrations down;

migrate-create:
	migrate create -ext sql -dir migrations -seq 'name';	


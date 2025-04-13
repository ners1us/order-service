CURDIR=$(shell pwd)
GENERATED_DIR=${CURDIR}/pkg/generated
PROTODIR=${GENERATED_DIR}/proto
MODULE=github.com/ners1us/order-service

.PHONY: generate-proto run stop clean-db rest-logs grpc-logs db-logs logs unit-test integration-test test help

.SILENT:

run:
	echo "Starting services..."
	docker compose up -d --build

stop:
	echo "Stopping services..."
	docker compose down

db-clean:
	echo "Removing database volume..."
	docker volume rm order-service_postgres-data

rest-logs:
	docker logs order-service-rest-app-1

rest-logs-follow:
	docker logs order-service-rest-app-1 -f

grpc-logs:
	docker logs order-service-grpc-app-1

grpc-logs-follow:
	docker logs order-service-grpc-app-1 -f

db-logs:
	docker logs order-service-postgres-1

db-logs-follow:
	docker logs order-service-postgres-1 -f

logs:
	echo "====== REST App Logs ======"
	docker logs order-service-rest-app-1 --tail 50
	echo "\n====== gRPC App Logs ======"
	docker logs order-service-grpc-app-1 --tail 50
	echo "\n====== Database Logs ======"
	docker logs order-service-postgres-1 --tail 50

unit-test:
	go test ./internal/services -v --cover

integration-test:
	go test ./internal/api/rest/... -v

test: unit-test integration-test
	echo "All tests completed"

generate-proto:
	echo "Generating protobuf code..."
	rm -rf ${PROTODIR}
	mkdir -p ${PROTODIR}
	protoc --proto_path=${CURDIR} \
		--go_out=module=${MODULE}:${CURDIR} \
		--go-grpc_out=module=${MODULE}:${CURDIR} \
		${CURDIR}/internal/api/grpc/proto/pvz.proto
	go mod tidy
	echo "Proto generation completed"

help:
	echo "Available commands:"
	echo "   make run                   - Start all containers"
	echo "   make stop                  - Stop all containers"
	echo "   make db-clean              - Remove database volume"
	echo "   make logs                  - Show all service logs"
	echo "   make rest-logs             - Show REST service logs"
	echo "   make rest-logs-follow      - Follow REST service logs"
	echo "   make grpc-logs             - Show gRPC service logs"
	echo "   make grpc-logs-follow      - Follow gRPC service logs"
	echo "   make db-logs               - Show database logs"
	echo "   make db-logs-follow        - Follow database logs"
	echo "   make unit-test             - Run unit tests with coverage"
	echo "   make integration-test      - Run integration tests"
	echo "   make test                  - Run all tests"
	echo "   make generate-proto        - Generate Go code from the proto file"

CURDIR=$(shell pwd)
GENERATED_DIR=${CURDIR}/pkg/generated
PROTODIR=${GENERATED_DIR}/proto

generate-proto:
	rm -rf ${PROTODIR}
	mkdir -p ${PROTODIR}
	protoc --proto_path=${CURDIR} \
	--go_out=module=github.com/ners1us/order-service:${CURDIR} \
	--go-grpc_out=module=github.com/ners1us/order-service:${CURDIR} \
	${CURDIR}/internal/api/grpc/pvz.proto
	go mod tidy
all: clean grpc build

clean:
	rm -rf grpc_api

build:
	go build

grpc:
	mkdir -p grpc_api/raft
	protoc -I proto --go_out=plugins:grpc_api/raft proto/raft-grpc.proto

run:
	./raft